package meilisearch

import "github.com/meilisearch/meilisearch-go"

func isStatusUnknown(taskStatus meilisearch.TaskStatus) bool {
	return taskStatus == meilisearch.TaskStatusUnknown
}

func isStatusEnqueued(taskStatus meilisearch.TaskStatus) bool {
	return taskStatus == meilisearch.TaskStatusEnqueued
}

func isStatusProcessing(taskStatus meilisearch.TaskStatus) bool {
	return taskStatus == meilisearch.TaskStatusProcessing
}

func isStatusSucceeded(taskStatus meilisearch.TaskStatus) bool {
	return taskStatus == meilisearch.TaskStatusSucceeded
}

func isStatusFailed(taskStatus meilisearch.TaskStatus) bool {
	return taskStatus == meilisearch.TaskStatusFailed
}

func isOk(taskStatus meilisearch.TaskStatus) bool {
	return isStatusEnqueued(taskStatus) || isStatusProcessing(taskStatus) || isStatusSucceeded(taskStatus)
}
