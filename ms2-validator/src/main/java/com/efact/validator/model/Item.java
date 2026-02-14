package com.efact.validator.model;

import lombok.Data;

@Data
public class Item {
    private String descripcion;
    private Double precioUnitario;
    private Integer cantidad;
    private Double precioTotal;
    private Double igvTotal;
}
