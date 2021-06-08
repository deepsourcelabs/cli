package api

import (
	"context"
	"fmt"

	"github.com/deepsourcelabs/graphql"
)

type IssuesListFileResponse struct {
	Repository struct {
		File struct {
			Issues struct {
				Edges []struct {
					Node struct {
						Path          string `json:"path"`
						Beginline     int    `json:"beginLine"`
						Endline       int    `json:"endLine"`
						Concreteissue struct {
							Analyzer struct {
								Shortcode string `json:"shortcode"`
							} `json:"analyzer"`
							Title     string `json:"title"`
							Shortcode string `json:"shortcode"`
						} `json:"concreteIssue"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"issues"`
		} `json:"file"`
	} `json:"repository"`
}

func GetIssuesForFile(client *DSClient, owner string, repoName string, provider string, filePath string, limit int) (*IssuesListFileResponse, error) {

	gq := client.gqlClient

	query := `
    query($name:String!, $owner:String!, $provider:VCSProviderChoices!, $path:String!, $limit:Int!){
        repository(name:$name, owner:$owner, provider:$provider){
            file(path:$path){
                issues(first:$limit){
                    edges{
                        node{
                            path
                            beginLine
                            endLine
                            concreteIssue{
                                analyzer {
                                    shortcode
                                }
                                title
                                shortcode
                            }
                        }
                    }
                }
            }
        }
    }`

	req := graphql.NewRequest(query)
	req.Var("name", repoName)
	req.Var("owner", owner)
	req.Var("provider", provider)
	req.Var("path", filePath)
	req.Var("limit", limit)

	// set header fields
	header := fmt.Sprintf("JWT %s", client.Token)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Add("Authorization", header)

	// define a Context for the request
	ctx := context.Background()

	// run it and capture the response
	// var graphqlResponse map[string]interface{}
	var respData IssuesListFileResponse
	if err := gq.Run(ctx, req, &respData); err != nil {
		return &respData, err
	}

	return &respData, nil

}
