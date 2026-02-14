package com.efact.validator.messaging;

import com.efact.validator.constants.MessageConstants;
import com.efact.validator.constants.QueueConstants;
import com.efact.validator.model.Documento;
import com.efact.validator.service.ISignatureService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.amqp.core.Message;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;

import java.util.HashMap;
import java.util.Map;

@Component
public class VerifyConsumer {

    private static final Logger logger = LoggerFactory.getLogger(VerifyConsumer.class);

    private final ISignatureService signatureService;
    private final RabbitTemplate rabbitTemplate;
    private final ObjectMapper objectMapper;

    public VerifyConsumer(
            ISignatureService signatureService,
            RabbitTemplate rabbitTemplate,
            ObjectMapper objectMapper) {
        this.signatureService = signatureService;
        this.rabbitTemplate = rabbitTemplate;
        this.objectMapper = objectMapper;
    }

    @RabbitListener(queues = QueueConstants.VERIFY_REQUEST_QUEUE)
    public void receiveVerifyRequest(Message message) {
        try {
            String correlationId = message.getMessageProperties().getCorrelationId();
            String replyTo = message.getMessageProperties().getReplyTo();

            logger.info("Solicitud de verificacion recibida - CorrelationId: {}", correlationId);

            String body = new String(message.getBody());
            Map<String, Object> request = objectMapper.readValue(body, Map.class);

            Map<String, Object> documentoMap = (Map<String, Object>) request.get("documento");
            String firma = (String) request.get("firma");

            Documento documento = objectMapper.convertValue(documentoMap, Documento.class);

            logger.info("Documento recibido para verificar: {}", objectMapper.writeValueAsString(documento));
            logger.info("Firma recibida: {}", firma);

            boolean valido = signatureService.verifySignature(documento, firma);

            Map<String, Object> response = new HashMap<>();
            response.put("valido", valido);
            if (valido) {
                response.put("mensaje", MessageConstants.FIRMA_VALIDA);
            } else {
                response.put("mensaje", MessageConstants.FIRMA_INVALIDA);
            }

            rabbitTemplate.convertAndSend(replyTo, response, msg -> {
                msg.getMessageProperties().setCorrelationId(correlationId);
                msg.getMessageProperties().setContentType("application/json");
                return msg;
            });

            logger.info("Respuesta de verificacion enviada - CorrelationId: {}, Valido: {}", correlationId, valido);

        } catch (Exception e) {
            logger.error("Error al procesar solicitud de verificacion", e);
        }
    }
}
