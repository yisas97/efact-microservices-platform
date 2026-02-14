package com.efact.validator.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import lombok.Data;
import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;
import org.springframework.data.mongodb.core.mapping.Field;

import java.util.List;

@Data
@Document(collection = "documents")
public class Documento {
    @Id
    @JsonIgnore
    private String id;

    @Field("idDocumento")
    private String idDocumento;

    @Field("uuid")
    private String uuid;

    @Field("rucEmisor")
    private String rucEmisor;

    @Field("rucReceptor")
    private String rucReceptor;

    @Field("fechaEmision")
    private String fechaEmision;

    @Field("montoTotalSinImpuestos")
    private Double montoTotalSinImpuestos;

    @Field("igvTotal")
    private Double igvTotal;

    @Field("montoTotal")
    private Double montoTotal;

    @Field("items")
    private List<Item> items;

    @Field("validacion")
    private Validacion validacion;
}
