package model

type JobType string

const (
	JobTypeWeibo = JobType("weibo")
)

type JobState string

const (
	JobStateWaiting = JobState("waiting")
	JobStateRunning = JobState("running")
	JobStateFailed  = JobState("failed")
	JobStateSucceed = JobState("succeed")
)
