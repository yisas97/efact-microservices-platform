package com.efact.validator.messaging;

import com.efact.validator.model.DocumentMessage;
import com.efact.validator.service.IDocumentProcessorService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;

@Component
public class DocumentConsumer {

    private static final Logger logger = LoggerFactory.getLogger(DocumentConsumer.class);

    private final IDocumentProcessorService documentProcessorService;

    public DocumentConsumer(IDocumentProcessorService documentProcessorService) {
        this.documentProcessorService = documentProcessorService;
    }

    @RabbitListener(queues = "${rabbitmq.queue.name}")
    public void receiveMessage(DocumentMessage message) {
        logger.info("Mensaje recibido de RabbitMQ: documentId={}, uuid={}",
            message.getDocumentId(), message.getUuid());

        try {
            documentProcessorService.processDocument(message.getDocumentId());
            logger.info("Mensaje procesado exitosamente");
        } catch (Exception e) {
            logger.error("Error al procesar el mensaje", e);
        }
    }
}
