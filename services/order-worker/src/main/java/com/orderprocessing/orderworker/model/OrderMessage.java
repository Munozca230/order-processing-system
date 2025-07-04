package com.orderprocessing.orderworker.model;

import java.util.List;

public record OrderMessage(String orderId, String customerId, List<ProductReference> products) {
}
