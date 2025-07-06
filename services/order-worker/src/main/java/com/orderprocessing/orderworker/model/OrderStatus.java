package com.orderprocessing.orderworker.model;

public enum OrderStatus {
    PROCESSING,  // Order in Kafka queue or being processed
    COMPLETED,   // Successfully processed and stored
    FAILED,      // Failed after max retries
    RETRYING     // Currently in retry queue
}