package com.efact.validator.service;

import com.efact.validator.exception.SignatureException;
import com.efact.validator.model.Documento;
import com.efact.validator.util.DocumentUtils;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.security.*;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

@Service
public class SignatureServiceImpl implements ISignatureService {

    private static final Logger logger = LoggerFactory.getLogger(SignatureServiceImpl.class);

    private final ObjectMapper objectMapper;

    @Value("${signature.private-key}")
    private String privateKeyBase64;

    @Value("${signature.public-key}")
    private String publicKeyBase64;

    private PrivateKey privateKey;
    private PublicKey publicKey;

    public SignatureServiceImpl(ObjectMapper objectMapper) {
        this.objectMapper = objectMapper;
    }

    @PostConstruct
    public void init() {
        loadKeys();
        logger.info("Servicio de firmas inicializado con claves RSA desde configuración");
    }

    private void loadKeys() {
        try {
            KeyFactory keyFactory = KeyFactory.getInstance("RSA");

            byte[] privateKeyBytes = Base64.getDecoder().decode(privateKeyBase64);
            PKCS8EncodedKeySpec privateKeySpec = new PKCS8EncodedKeySpec(privateKeyBytes);
            privateKey = keyFactory.generatePrivate(privateKeySpec);

            byte[] publicKeyBytes = Base64.getDecoder().decode(publicKeyBase64);
            X509EncodedKeySpec publicKeySpec = new X509EncodedKeySpec(publicKeyBytes);
            publicKey = keyFactory.generatePublic(publicKeySpec);

            logger.info("Claves RSA cargadas exitosamente desde configuración");
        } catch (Exception e) {
            logger.error("Error al cargar las claves RSA desde configuración", e);
            throw new SignatureException("Fallo al cargar las claves RSA", e);
        }
    }

    @Override
    public String signDocument(Documento document) {
        try {
            Documento docCopy = DocumentUtils.cloneDocument(document, objectMapper);
            docCopy.setValidacion(null);

            String documentJson = objectMapper.writeValueAsString(docCopy);
            logger.info("JSON generado al firmar documento {}: {}", document.getIdDocumento(), documentJson);

            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(documentJson.getBytes());

            Signature signature = Signature.getInstance("SHA256withRSA");
            signature.initSign(privateKey);
            signature.update(hash);
            byte[] signedHash = signature.sign();

            String base64Signature = Base64.getEncoder().encodeToString(signedHash);

            logger.info("Documento firmado exitosamente: {}", document.getIdDocumento());
            return base64Signature;
        } catch (Exception e) {
            logger.error("Error al firmar el documento", e);
            throw new SignatureException("Fallo al firmar el documento", e);
        }
    }

    @Override
    public boolean verifySignature(Documento document, String signature) {
        try {
            Documento docCopy = DocumentUtils.cloneDocument(document, objectMapper);
            docCopy.setValidacion(null);

            String documentJson = objectMapper.writeValueAsString(docCopy);
            logger.info("JSON generado para verificar: {}", documentJson);

            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(documentJson.getBytes());

            byte[] signatureBytes = Base64.getDecoder().decode(signature);

            Signature sig = Signature.getInstance("SHA256withRSA");
            sig.initVerify(publicKey);
            sig.update(hash);

            return sig.verify(signatureBytes);
        } catch (Exception e) {
            logger.error("Error al verificar la firma", e);
            return false;
        }
    }

    @Override
    public PublicKey getPublicKey() {
        return publicKey;
    }

    @Override
    public String getPublicKeyBase64() {
        return publicKeyBase64;
    }
}
