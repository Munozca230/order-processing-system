package com.example.orderworker.model;

import java.util.List;

public record EnrichedOrder(OrderMessage order, CustomerDetails customer, List<ProductDetails> products) {}
