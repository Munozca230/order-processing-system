package com.example.orderworker.model;

import java.time.LocalDateTime;

public class FailedMessage {
    private String messageId;
    private String content;
    private String errorMessage;
    private int retryCount;
    private LocalDateTime lastRetryAt;
    private LocalDateTime nextRetryAt;
    private String failureReason;

    public FailedMessage() {}

    public FailedMessage(String messageId, String content, String errorMessage) {
        this.messageId = messageId;
        this.content = content;
        this.errorMessage = errorMessage;
        this.retryCount = 0;
        this.lastRetryAt = LocalDateTime.now();
        this.nextRetryAt = calculateNextRetry(0);
        this.failureReason = "API_CALL_FAILED";
    }

    private LocalDateTime calculateNextRetry(int currentRetryCount) {
        // Exponential backoff: 1s, 2s, 4s, 8s, 16s
        long delaySeconds = (long) Math.pow(2, currentRetryCount);
        return LocalDateTime.now().plusSeconds(delaySeconds);
    }

    public void incrementRetry() {
        this.retryCount++;
        this.lastRetryAt = LocalDateTime.now();
        this.nextRetryAt = calculateNextRetry(this.retryCount);
    }

    // Getters and Setters
    public String getMessageId() {
        return messageId;
    }

    public void setMessageId(String messageId) {
        this.messageId = messageId;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getErrorMessage() {
        return errorMessage;
    }

    public void setErrorMessage(String errorMessage) {
        this.errorMessage = errorMessage;
    }

    public int getRetryCount() {
        return retryCount;
    }

    public void setRetryCount(int retryCount) {
        this.retryCount = retryCount;
    }

    public LocalDateTime getLastRetryAt() {
        return lastRetryAt;
    }

    public void setLastRetryAt(LocalDateTime lastRetryAt) {
        this.lastRetryAt = lastRetryAt;
    }

    public LocalDateTime getNextRetryAt() {
        return nextRetryAt;
    }

    public void setNextRetryAt(LocalDateTime nextRetryAt) {
        this.nextRetryAt = nextRetryAt;
    }

    public String getFailureReason() {
        return failureReason;
    }

    public void setFailureReason(String failureReason) {
        this.failureReason = failureReason;
    }
}