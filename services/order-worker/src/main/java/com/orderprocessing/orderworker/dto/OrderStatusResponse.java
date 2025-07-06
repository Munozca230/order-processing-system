package com.orderprocessing.orderworker.dto;

import com.orderprocessing.orderworker.model.OrderStatus;
import java.time.LocalDateTime;

public record OrderStatusResponse(
    String orderId,
    String status,
    LocalDateTime createdAt,
    LocalDateTime processedAt,
    LocalDateTime updatedAt,
    String failureReason,
    int retryCount,
    LocalDateTime nextRetryAt
) {
    public static OrderStatusResponse processing(String orderId) {
        return new OrderStatusResponse(
            orderId,
            OrderStatus.PROCESSING.name(),
            LocalDateTime.now(),
            null,
            LocalDateTime.now(),
            null,
            0,
            null
        );
    }
}