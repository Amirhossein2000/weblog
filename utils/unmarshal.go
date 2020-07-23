package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"weblog/schema"
)

func UnmarshalArticle(body io.ReadCloser) (*schema.Article, error) {
	byteBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	article := &schema.Article{}
	err = json.Unmarshal(byteBody, article)
	if err != nil {
		return nil, err
	}

	return article, nil
}

func UnmarshalComment(body io.ReadCloser) (*schema.Comment, error) {
	byteBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	comment := &schema.Comment{}
	err = json.Unmarshal(byteBody, comment)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func UnmarshalUser(body io.ReadCloser) (*schema.User, error) {
	byteBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	user := &schema.User{}
	err = json.Unmarshal(byteBody, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
