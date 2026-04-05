package main

import (
	"github.com/99designs/gqlgen/graphql"
	account "github.com/wignn/micro-3/account/client"
	auth "github.com/wignn/micro-3/auth/client"
	novel "github.com/wignn/micro-3/novel/client"
	readinglist "github.com/wignn/micro-3/readinglist/client"
	review "github.com/wignn/micro-3/review/client"
)

type GraphQLServer struct {
	accountClient     *account.AccountClient
	novelClient       *novel.NovelClient
	readinglistClient *readinglist.ReadingListClient
	authClient        *auth.AuthClient
	reviewClient      *review.ReviewClient
}

func NewGraphQLServer(accountUrl, novelUrl, readinglistUrl, reviewUrl, authUrl string) (*GraphQLServer, error) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	novelClient, err := novel.NewClient(novelUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}

	readinglistClient, err := readinglist.NewClient(readinglistUrl)
	if err != nil {
		novelClient.Close()
		return nil, err
	}

	reviewClient, err := review.NewClient(reviewUrl)
	if err != nil {
		readinglistClient.Close()
		return nil, err
	}

	authClient, err := auth.NewClient(authUrl)
	if err != nil {
		reviewClient.Close()
		return nil, err
	}

	return &GraphQLServer{
		accountClient,
		novelClient,
		readinglistClient,
		authClient,
		reviewClient,
	}, nil
}

func (s *GraphQLServer) Mutation() MutationResolver {
	return &mutationResolver{server: s}
}

func (s *GraphQLServer) Query() QueryResolver {
	return &queryResolver{server: s}
}

func (s *GraphQLServer) Account() AccountResolver {
	return &accountResolver{server: s}
}

func (s *GraphQLServer) ToExecutableSchema() (graphql.ExecutableSchema, error) {
	return NewExecutableSchema(Config{
		Resolvers: s,
	}), nil
}
