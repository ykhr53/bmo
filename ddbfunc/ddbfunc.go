package ddbfunc

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// GetVote gets vote value from DynamoDB
func GetVote(ddb *dynamodb.DynamoDB, key string) (string, error) {
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
		return "err", err
	}
	if len(resp.Item) == 0 {
		return "unvoted", nil
	}
	return *resp.Item["votes"].N, nil
}

// SetVote sets vote value to DynamoDB
func SetVote(ddb *dynamodb.DynamoDB, key string, val string) {
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
				N: aws.String(val),
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

// GetWord gets word value from DynamoDB
func GetWord(ddb *dynamodb.DynamoDB, key string) (string, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String("bmo"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(key),
			},
		},
		AttributesToGet: []*string{
			aws.String("description"),
		},
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: aws.String("NONE"),
	}

	resp, err := ddb.GetItem(params)
	if err != nil {
		fmt.Println(err.Error())
		return "err", err
	}
	if len(resp.Item) == 0 {
		return "unknown", nil
	}
	return *resp.Item["description"].S, nil
}

// SetWord sets word value to DynamoDB
func SetWord(ddb *dynamodb.DynamoDB, key string, val string) {
	param := &dynamodb.UpdateItemInput{
		TableName: aws.String("bmo"),
		Key: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(key),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#description": aws.String("description"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":desc_val": {
				S: aws.String(val),
			},
		},
		UpdateExpression:            aws.String("set #description = :desc_val"),
		ReturnConsumedCapacity:      aws.String("NONE"),
		ReturnItemCollectionMetrics: aws.String("NONE"),
		ReturnValues:                aws.String("NONE"),
	}

	_, err := ddb.UpdateItem(param)
	if err != nil {
		fmt.Println(err.Error())
	}
}
