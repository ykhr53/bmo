package ddbfunc

import (
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetVal gets value from DynamoDB
func GetVal(ddb *dynamodb.DynamoDB, key string) (int, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("bmo"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(key),
			},
		},
		AttributesToGet: []*string{
			aws.String("votes"),
		},
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: aws.String("NONE"),
	}

	resp, err := ddb.GetItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return -1, nil
	}
	if len(resp.Item) == 0 {
		return -1, nil
	}
	return strconv.Atoi(*resp.Item["votes"].N)
}

// SetVal sets value from DynamoDB
func SetVal(ddb *dynamodb.DynamoDB, key string, n string) {
	param := &dynamodb.UpdateItemInput{
		TableName: aws.String("bmo"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(key),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#votes": aws.String("votes"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":vote_val": {
				N: aws.String(n),
			},
		},
		UpdateExpression:            aws.String("set #votes = :vote_val"),
		ReturnConsumedCapacity:      aws.String("NONE"),
		ReturnItemCollectionMetrics: aws.String("NONE"),
		ReturnValues:                aws.String("NONE"),
	}

	_, err := ddb.UpdateItem(param)
	if err != nil {
		fmt.Println(err.Error())
	}
}
