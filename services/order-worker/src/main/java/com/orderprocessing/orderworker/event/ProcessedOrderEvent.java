package com.orderprocessing.orderworker.event;

import com.orderprocessing.orderworker.model.EnrichedOrder;
import org.springframework.context.ApplicationEvent;

public class ProcessedOrderEvent extends ApplicationEvent {
    private final EnrichedOrder order;
    public ProcessedOrderEvent(Object source, EnrichedOrder order) {
        super(source);
        this.order = order;
    }
    public EnrichedOrder getOrder() { return order; }
}
