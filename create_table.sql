CREATE TABLE public.estatistica
(
    cpf character varying(14) NOT NULL,
    "cpfvalido" boolean NOT NULL,
    "dataultimacompra" date,
    incompleto boolean NOT NULL,
    "lojamaisfrequente" character varying(14),
    "lojaultimacompra" character varying(14),
    private boolean NOT NULL,
    "ticketmedio" numeric(10, 2),
    "ticketultimacompra" numeric(10, 2)
    --,PRIMARY KEY (cpf)
);

ALTER TABLE public.estatistica
    OWNER to postgres;
