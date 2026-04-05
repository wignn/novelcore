#!/usr/bin/env python3
"""Smoke test untuk semua endpoint GraphQL di schema.

Jalankan ini hanya ke environment development/test karena script akan create,
update, dan delete data.
"""

from __future__ import annotations

import argparse
import json
import os
import sys
import time
import uuid
from dataclasses import dataclass
from typing import Any, Dict, List, Optional
from urllib import error, request


class GraphQLRequestError(Exception):
    pass


@dataclass
class TestResult:
    name: str
    ok: bool
    skipped: bool
    message: str
    elapsed_ms: int


class GraphQLClient:
    def __init__(self, url: str, timeout: float = 10.0) -> None:
        self.url = url
        self.timeout = timeout

    def execute(self, query: str, variables: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        payload = {"query": query, "variables": variables or {}}
        data = json.dumps(payload).encode("utf-8")
        req = request.Request(
            self.url,
            data=data,
            headers={"Content-Type": "application/json"},
            method="POST",
        )

        try:
            with request.urlopen(req, timeout=self.timeout) as resp:
                raw = resp.read().decode("utf-8")
        except error.HTTPError as exc:
            body = exc.read().decode("utf-8", errors="replace")
            raise GraphQLRequestError(f"HTTP {exc.code}: {body}") from exc
        except error.URLError as exc:
            raise GraphQLRequestError(f"Tidak bisa konek ke {self.url}: {exc}") from exc

        try:
            parsed = json.loads(raw)
        except json.JSONDecodeError as exc:
            raise GraphQLRequestError(f"Response bukan JSON valid: {raw[:300]}") from exc

        errs = parsed.get("errors")
        if errs:
            msg = "; ".join(str(e.get("message", e)) for e in errs)
            raise GraphQLRequestError(msg)

        return parsed.get("data", {})


def require(value: Any, label: str) -> Any:
    if value is None:
        raise AssertionError(f"Nilai kosong untuk {label}")
    return value


def run_test(name: str, fn) -> TestResult:
    start = time.perf_counter()
    try:
        fn()
        elapsed = int((time.perf_counter() - start) * 1000)
        return TestResult(name=name, ok=True, skipped=False, message="OK", elapsed_ms=elapsed)
    except Exception as exc:  # noqa: BLE001
        elapsed = int((time.perf_counter() - start) * 1000)
        return TestResult(name=name, ok=False, skipped=False, message=str(exc), elapsed_ms=elapsed)


def make_skip(name: str, reason: str) -> TestResult:
    return TestResult(name=name, ok=True, skipped=True, message=reason, elapsed_ms=0)


def mask_token(token: str, head: int = 10, tail: int = 6) -> str:
    if len(token) <= head + tail:
        return token
    return f"{token[:head]}...{token[-tail:]}"


def main() -> int:
    parser = argparse.ArgumentParser(description="Test semua endpoint GraphQL")
    parser.add_argument(
        "--url",
        default=os.getenv("GQL_URL", "http://localhost:8000/graphql"),
        help="URL GraphQL endpoint (default: %(default)s)",
    )
    parser.add_argument(
        "--timeout",
        type=float,
        default=10.0,
        help="Timeout request dalam detik (default: %(default)s)",
    )
    parser.add_argument(
        "--keep-data",
        action="store_true",
        help="Jangan hapus data test di akhir",
    )
    args = parser.parse_args()

    client = GraphQLClient(args.url, args.timeout)
    seed = uuid.uuid4().hex[:8]

    state: Dict[str, Any] = {
        "account_id": None,
        "author_id": None,
        "group_id": None,
        "novel_id": None,
        "chapter_id": None,
        "review_id": None,
        "readinglist_id": None,
        "refresh_token": None,
        "account_refresh_token": None,
        "account_access_token": None,
        "genre_ids": [],
        "tag_ids": [],
        "email": f"gqltest_{seed}@example.com",
        "password": "Pass1234!",
    }

    results: List[TestResult] = []

    def add(name: str, fn) -> None:
        res = run_test(name, fn)
        results.append(res)
        if res.skipped:
            status = "SKIP"
        else:
            status = "PASS" if res.ok else "FAIL"
        print(f"[{status}] {res.name} ({res.elapsed_ms} ms) - {res.message}")

    def add_with_deps(name: str, deps: List[str], fn) -> None:
        missing = [d for d in deps if not state.get(d)]
        if missing:
            res = make_skip(name, f"Prasyarat belum ada: {', '.join(missing)}")
            results.append(res)
            print(f"[SKIP] {res.name} ({res.elapsed_ms} ms) - {res.message}")
            return
        add(name, fn)

    add(
        "Query genres",
        lambda: state.update(
            {
                "genre_ids": [
                    g["id"]
                    for g in client.execute("query { genres { id name } }")["genres"]
                ]
            }
        ),
    )

    add(
        "Query tags",
        lambda: state.update(
            {
                "tag_ids": [
                    t["id"] for t in client.execute("query { tags { id name } }")["tags"]
                ]
            }
        ),
    )

    add(
        "Mutation createAccount",
        lambda: state.update(
            {
                "account_id": require(
                    client.execute(
                        """
                        mutation($in: AccountInput!) {
                          createAccount(account: $in) { id email }
                        }
                        """,
                        {
                            "in": {
                                "name": f"GQL Tester {seed}",
                                "email": state["email"],
                                "password": state["password"],
                            }
                        },
                    )["createAccount"]["id"],
                    "createAccount.id",
                )
            }
        ),
    )

    def login_and_capture_tokens() -> None:
        data = client.execute(
            """
            mutation($in: LoginInput!) {
              login(account: $in) {
                id
                backendToken {
                  accessToken
                  refreshToken
                }
              }
            }
            """,
            {"in": {"email": state["email"], "password": state["password"]}},
        )["login"]

        login_id = require(data["id"], "login.id")
        if state["account_id"] and login_id != state["account_id"]:
            raise AssertionError("login.id tidak sama dengan account_id hasil createAccount")

        refresh_token = require(data["backendToken"]["refreshToken"], "login.backendToken.refreshToken")
        access_token = require(data["backendToken"]["accessToken"], "login.backendToken.accessToken")

        state["account_refresh_token"] = refresh_token
        state["refresh_token"] = refresh_token
        state["account_access_token"] = access_token

    add_with_deps(
        "Mutation login (ambil refresh token akun)",
        ["account_id"],
        login_and_capture_tokens,
    )

    add_with_deps(
        "Mutation refreshToken",
        ["refresh_token"],
        lambda: require(
            client.execute(
                """
                mutation($refreshToken: String!) {
                  refreshToken(refreshToken: $refreshToken) { accessToken }
                }
                """,
                {"refreshToken": state["refresh_token"]},
            )["refreshToken"]["accessToken"],
            "refreshToken.accessToken",
        ),
    )

    add_with_deps(
        "Mutation editAccount",
        ["account_id"],
        lambda: require(
            client.execute(
                """
                mutation($id: String!, $in: EditAccountInput!) {
                  editAccount(id: $id, account: $in) { id name bio }
                }
                """,
                {
                    "id": state["account_id"],
                    "in": {"name": f"GQL Tester Updated {seed}", "bio": "bio test"},
                },
            )["editAccount"]["id"],
            "editAccount.id",
        ),
    )

    add_with_deps(
        "Query accounts (id)",
        ["account_id"],
        lambda: require(
            client.execute(
                "query($id: String) { accounts(id: $id) { id email } }",
                {"id": state["account_id"]},
            )["accounts"],
            "accounts",
        ),
    )

    add(
        "Mutation createAuthor",
        lambda: state.update(
            {
                "author_id": require(
                    client.execute(
                        """
                        mutation($in: AuthorInput!) {
                          createAuthor(author: $in) { id name }
                        }
                        """,
                        {"in": {"name": f"Author {seed}", "bio": "author bio"}},
                    )["createAuthor"]["id"],
                    "createAuthor.id",
                )
            }
        ),
    )

    add_with_deps(
        "Query authors (id)",
        ["author_id"],
        lambda: require(
            client.execute(
                "query($id: String) { authors(id: $id) { id name } }",
                {"id": state["author_id"]},
            )["authors"],
            "authors",
        ),
    )

    add(
        "Mutation createTranslationGroup",
        lambda: state.update(
            {
                "group_id": require(
                    client.execute(
                        """
                        mutation($in: TranslationGroupInput!) {
                          createTranslationGroup(group: $in) { id name }
                        }
                        """,
                        {
                            "in": {
                                "name": f"Group {seed}",
                                "websiteUrl": "https://example.com",
                                "description": "group test",
                            }
                        },
                    )["createTranslationGroup"]["id"],
                    "createTranslationGroup.id",
                )
            }
        ),
    )

    add(
        "Query translationGroups",
        lambda: require(
            client.execute("query { translationGroups { id name } }")["translationGroups"],
            "translationGroups",
        ),
    )

    add_with_deps(
        "Mutation createNovel",
        ["author_id"],
        lambda: state.update(
            {
                "novel_id": require(
                    client.execute(
                        """
                        mutation($in: NovelInput!) {
                          createNovel(novel: $in) { id title }
                        }
                        """,
                        {
                            "in": {
                                "title": f"Novel {seed}",
                                "description": "novel test",
                                "authorId": state["author_id"],
                                "status": "ongoing",
                                "novelType": "web_novel",
                                "countryOfOrigin": "ID",
                                "yearPublished": 2026,
                                "genreIds": state["genre_ids"][:2],
                                "tagIds": state["tag_ids"][:2],
                            }
                        },
                    )["createNovel"]["id"],
                    "createNovel.id",
                )
            }
        ),
    )

    add_with_deps(
        "Mutation updateNovel",
        ["novel_id", "author_id"],
        lambda: require(
            client.execute(
                """
                mutation($id: String!, $in: NovelInput!) {
                  updateNovel(id: $id, novel: $in) { id title }
                }
                """,
                {
                    "id": state["novel_id"],
                    "in": {
                        "title": f"Novel Updated {seed}",
                        "description": "novel test updated",
                        "authorId": state["author_id"],
                        "status": "ongoing",
                        "novelType": "web_novel",
                        "countryOfOrigin": "ID",
                        "yearPublished": 2026,
                        "genreIds": state["genre_ids"][:2],
                        "tagIds": state["tag_ids"][:2],
                    },
                },
            )["updateNovel"]["id"],
            "updateNovel.id",
        ),
    )

    add_with_deps(
        "Query novels (id)",
        ["novel_id"],
        lambda: require(
            client.execute(
                "query($id: String) { novels(id: $id) { id title } }",
                {"id": state["novel_id"]},
            )["novels"],
            "novels",
        ),
    )

    add_with_deps(
        "Mutation createChapter",
        ["novel_id", "group_id"],
        lambda: state.update(
            {
                "chapter_id": require(
                    client.execute(
                        """
                        mutation($in: ChapterInput!) {
                          createChapter(chapter: $in) { id chapterNumber }
                        }
                        """,
                        {
                            "in": {
                                "novelId": state["novel_id"],
                                "chapterNumber": 1.0,
                                "title": "Chapter 1",
                                "translatorGroupId": state["group_id"],
                                "sourceUrl": "https://example.com/ch1",
                            }
                        },
                    )["createChapter"]["id"],
                    "createChapter.id",
                )
            }
        ),
    )

    add_with_deps(
        "Mutation updateChapter",
        ["chapter_id", "novel_id", "group_id"],
        lambda: require(
            client.execute(
                """
                mutation($id: String!, $in: ChapterInput!) {
                  updateChapter(id: $id, chapter: $in) { id chapterNumber }
                }
                """,
                {
                    "id": state["chapter_id"],
                    "in": {
                        "novelId": state["novel_id"],
                        "chapterNumber": 2.0,
                        "title": "Chapter 2",
                        "translatorGroupId": state["group_id"],
                        "sourceUrl": "https://example.com/ch2",
                    },
                },
            )["updateChapter"]["id"],
            "updateChapter.id",
        ),
    )

    add_with_deps(
        "Query chapters",
        ["novel_id"],
        lambda: require(
            client.execute(
                "query($novelId: String!) { chapters(novelId: $novelId) { id chapterNumber } }",
                {"novelId": state["novel_id"]},
            )["chapters"],
            "chapters",
        ),
    )

    add_with_deps(
        "Query chapter",
        ["chapter_id"],
        lambda: require(
            client.execute(
                "query($id: String!) { chapter(id: $id) { id title } }",
                {"id": state["chapter_id"]},
            )["chapter"]["id"],
            "chapter.id",
        ),
    )

    add_with_deps(
        "Mutation addToReadingList",
        ["account_id", "novel_id"],
        lambda: state.update(
            {
                "readinglist_id": require(
                    client.execute(
                        """
                        mutation($accountId: String!, $entry: ReadingListInput!) {
                          addToReadingList(accountId: $accountId, entry: $entry) { id status }
                        }
                        """,
                        {
                            "accountId": state["account_id"],
                            "entry": {
                                "novelId": state["novel_id"],
                                "status": "reading",
                                "currentChapter": 2.0,
                                "rating": 8,
                                "notes": "test notes",
                                "isFavorite": True,
                            },
                        },
                    )["addToReadingList"]["id"],
                    "addToReadingList.id",
                )
            }
        ),
    )

    add_with_deps(
        "Mutation updateReadingList",
        ["readinglist_id", "novel_id"],
        lambda: require(
            client.execute(
                """
                mutation($id: String!, $entry: ReadingListInput!) {
                  updateReadingList(id: $id, entry: $entry) { id status }
                }
                """,
                {
                    "id": state["readinglist_id"],
                    "entry": {
                        "novelId": state["novel_id"],
                        "status": "completed",
                        "currentChapter": 2.0,
                        "rating": 9,
                        "notes": "updated notes",
                        "isFavorite": False,
                    },
                },
            )["updateReadingList"]["id"],
            "updateReadingList.id",
        ),
    )

    add_with_deps(
        "Query readingList",
        ["account_id"],
        lambda: require(
            client.execute(
                "query($accountId: String!) { readingList(accountId: $accountId) { id status } }",
                {"accountId": state["account_id"]},
            )["readingList"],
            "readingList",
        ),
    )

    add_with_deps(
        "Mutation createReview",
        ["novel_id", "account_id"],
        lambda: state.update(
            {
                "review_id": require(
                    client.execute(
                        """
                        mutation($in: ReviewInput!) {
                          createReview(review: $in) { id rating }
                        }
                        """,
                        {
                            "in": {
                                "novelId": state["novel_id"],
                                "accountId": state["account_id"],
                                "rating": 5,
                                "title": "Great",
                                "content": "Nice novel",
                                "isSpoiler": False,
                            }
                        },
                    )["createReview"]["id"],
                    "createReview.id",
                )
            }
        ),
    )

    add_with_deps(
        "Query reviews",
        ["novel_id"],
        lambda: require(
            client.execute(
                "query($novelId: String!) { reviews(novelId: $novelId) { id rating } }",
                {"novelId": state["novel_id"]},
            )["reviews"],
            "reviews",
        ),
    )

    add_with_deps(
        "Mutation incrementViewCount",
        ["novel_id"],
        lambda: require(
            client.execute(
                "mutation($novelId: String!) { incrementViewCount(novelId: $novelId) }",
                {"novelId": state["novel_id"]},
            )["incrementViewCount"],
            "incrementViewCount",
        ),
    )

    add(
        "Query novelRanking",
        lambda: require(
            client.execute(
                "query { novelRanking(period: WEEKLY, sortBy: VIEWS) { rank score } }"
            )["novelRanking"],
            "novelRanking",
        ),
    )

    # Delete endpoints juga dites sebagai bagian flow.
    add_with_deps(
        "Mutation deleteReview",
        ["review_id"],
        lambda: require(
            client.execute(
                "mutation($id: String!) { deleteReview(id: $id) { success deletedId } }",
                {"id": state["review_id"]},
            )["deleteReview"]["success"],
            "deleteReview.success",
        ),
    )
    if any(r.name == "Mutation deleteReview" and not r.skipped and r.ok for r in results):
        state["review_id"] = None

    add_with_deps(
        "Mutation removeFromReadingList",
        ["readinglist_id"],
        lambda: require(
            client.execute(
                "mutation($id: String!) { removeFromReadingList(id: $id) { success deletedId } }",
                {"id": state["readinglist_id"]},
            )["removeFromReadingList"]["success"],
            "removeFromReadingList.success",
        ),
    )
    if any(r.name == "Mutation removeFromReadingList" and not r.skipped and r.ok for r in results):
        state["readinglist_id"] = None

    add_with_deps(
        "Mutation deleteChapter",
        ["chapter_id"],
        lambda: require(
            client.execute(
                "mutation($id: String!) { deleteChapter(id: $id) { success deletedId } }",
                {"id": state["chapter_id"]},
            )["deleteChapter"]["success"],
            "deleteChapter.success",
        ),
    )
    if any(r.name == "Mutation deleteChapter" and not r.skipped and r.ok for r in results):
        state["chapter_id"] = None

    add_with_deps(
        "Mutation deleteNovel",
        ["novel_id"],
        lambda: require(
            client.execute(
                "mutation($id: String!) { deleteNovel(id: $id) { success deletedId } }",
                {"id": state["novel_id"]},
            )["deleteNovel"]["success"],
            "deleteNovel.success",
        ),
    )
    if any(r.name == "Mutation deleteNovel" and not r.skipped and r.ok for r in results):
        state["novel_id"] = None

    add_with_deps(
        "Mutation deleteAccount",
        ["account_id"],
        lambda: require(
            client.execute(
                "mutation($id: String!) { deleteAccount(id: $id) { success deletedId } }",
                {"id": state["account_id"]},
            )["deleteAccount"]["success"],
            "deleteAccount.success",
        ),
    )
    if any(r.name == "Mutation deleteAccount" and not r.skipped and r.ok for r in results):
        state["account_id"] = None

    if not args.keep_data:
        # Cleanup defensive kalau ada test yang gagal di tengah.
        def safe_cleanup(op_name: str, query: str, variables: Dict[str, Any]) -> None:
            if not variables:
                return
            try:
                client.execute(query, variables)
                print(f"[CLEANUP] {op_name}: OK")
            except Exception as exc:  # noqa: BLE001
                print(f"[CLEANUP] {op_name}: {exc}")

        if state.get("review_id"):
            safe_cleanup(
                "deleteReview",
                "mutation($id: String!) { deleteReview(id: $id) { success } }",
                {"id": state["review_id"]},
            )
        if state.get("readinglist_id"):
            safe_cleanup(
                "removeFromReadingList",
                "mutation($id: String!) { removeFromReadingList(id: $id) { success } }",
                {"id": state["readinglist_id"]},
            )
        if state.get("chapter_id"):
            safe_cleanup(
                "deleteChapter",
                "mutation($id: String!) { deleteChapter(id: $id) { success } }",
                {"id": state["chapter_id"]},
            )
        if state.get("novel_id"):
            safe_cleanup(
                "deleteNovel",
                "mutation($id: String!) { deleteNovel(id: $id) { success } }",
                {"id": state["novel_id"]},
            )
        if state.get("account_id"):
            safe_cleanup(
                "deleteAccount",
                "mutation($id: String!) { deleteAccount(id: $id) { success } }",
                {"id": state["account_id"]},
            )

    total = len(results)
    failed = len([r for r in results if not r.ok])
    skipped = len([r for r in results if r.skipped])
    passed = total - failed - skipped

    print("\n===== RINGKASAN =====")
    print(f"Total: {total}")
    print(f"Lolos: {passed}")
    print(f"Skip: {skipped}")
    print(f"Gagal: {failed}")
    if state.get("account_refresh_token"):
        print(f"Refresh token akun: {mask_token(state['account_refresh_token'])}")

    return 0 if failed == 0 else 1


if __name__ == "__main__":
    sys.exit(main())
