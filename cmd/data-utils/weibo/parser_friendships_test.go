package weibo

import "testing"

func TestParserFriendships(t *testing.T) {
	parser := NewFriendshipParser()
	var metaData MetaData
	metaData.Url = "https://weibo.com/ajax/friendships/friends?relate=fans&page=1&uid=7362070198&type=all&newFollowerCount=0"
	jobList, err := parser.Parse(metaData)
	if err != nil {
		t.Fatal(err)
	}
	for _, job := range jobList {
		t.Log(job.MetaData)
	}
}
