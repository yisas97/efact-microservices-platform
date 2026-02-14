package com.efact.validator.service;

import com.efact.validator.model.Documento;

import java.security.PublicKey;

public interface ISignatureService {

    String signDocument(Documento document);

    boolean verifySignature(Documento document, String signature);

    PublicKey getPublicKey();

    String getPublicKeyBase64();
}
