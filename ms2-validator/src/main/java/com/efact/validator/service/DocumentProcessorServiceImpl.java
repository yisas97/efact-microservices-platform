package com.efact.validator.service;

import com.efact.validator.constants.MessageConstants;
import com.efact.validator.constants.ValidationConstants;
import com.efact.validator.model.Documento;
import com.efact.validator.model.Validacion;
import com.efact.validator.repository.DocumentRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.time.Instant;
import java.util.Optional;

@Service
public class DocumentProcessorServiceImpl implements IDocumentProcessorService {

    private static final Logger logger = LoggerFactory.getLogger(DocumentProcessorServiceImpl.class);

    private final DocumentRepository documentRepository;
    private final IValidationService validationService;
    private final ISignatureService signatureService;

    public DocumentProcessorServiceImpl(
            DocumentRepository documentRepository,
            IValidationService validationService,
            ISignatureService signatureService) {
        this.documentRepository = documentRepository;
        this.validationService = validationService;
        this.signatureService = signatureService;
    }

    @Override
    public void processDocument(String documentId) {
        logger.info("Procesando documento: {}", documentId);

        Optional<Documento> optionalDoc = documentRepository.findByIdDocumento(documentId);

        if (optionalDoc.isEmpty()) {
            logger.error("{}: {}", MessageConstants.DOCUMENTO_NO_ENCONTRADO, documentId);
            return;
        }

        Documento document = optionalDoc.get();

        boolean isValid = validationService.validateDocument(document);

        Validacion validacion = new Validacion();
        validacion.setFechaValidacion(Instant.now().toString());

        if (isValid) {
            String signature = signatureService.signDocument(document);
            validacion.setFirma(signature);
            validacion.setEstado(ValidationConstants.ESTADO_VALIDO);
            logger.info("Documento {} validado y firmado exitosamente", documentId);
        } else {
            validacion.setEstado(ValidationConstants.ESTADO_INVALIDO);
            logger.warn("Documento {} marcado como inválido", documentId);
        }

        document.setValidacion(validacion);
        documentRepository.save(document);

        logger.info("Documento {} procesado y actualizado en la base de datos", documentId);
    }
}
