--
-- PostgreSQL database dump
--

\restrict eVfo24ztCrjq5oR5h4R0FuvgCXc2hThwtUD7gGMwP4d5ltbVgmbhcQoZ9VbVp0D

-- Dumped from database version 15.14
-- Dumped by pg_dump version 15.14

-- Started on 2025-09-17 17:41:21

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 214 (class 1259 OID 16400)
-- Name: books; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.books (
    id integer NOT NULL,
    title text NOT NULL,
    author text NOT NULL,
    pages integer NOT NULL,
    owner_id integer
);


ALTER TABLE public.books OWNER TO postgres;

--
-- TOC entry 215 (class 1259 OID 16405)
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.books_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.books_id_seq OWNER TO postgres;

--
-- TOC entry 3337 (class 0 OID 0)
-- Dependencies: 215
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- TOC entry 217 (class 1259 OID 16587)
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    role character varying(20)
);


ALTER TABLE public.users OWNER TO postgres;

--
-- TOC entry 216 (class 1259 OID 16586)
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO postgres;

--
-- TOC entry 3338 (class 0 OID 0)
-- Dependencies: 216
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- TOC entry 3178 (class 2604 OID 16406)
-- Name: books id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- TOC entry 3179 (class 2604 OID 16590)
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- TOC entry 3328 (class 0 OID 16400)
-- Dependencies: 214
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.books (id, title, author, pages, owner_id) FROM stdin;
5	Naturalis Historia	Gaius Plinius Secundus	418	1
6	Tusculanae Disputationes	Marcus Tullius Cicero	361	1
7	De Rerum Natura	Titus Lucretius Carus	248	1
8	Satyricon	Gaius Petronius Arbiter	157	1
9	Epistulae Morales	Lucius Annaeus Seneca	216	2
10	De Architectura	Marcus Vitruvius Pollio	168	2
11	Historiarum Libri	Publius Cornelius Tacitus	386	1
12	De Vita Caesarum	Gaius Suetonius Tranquillus	123	2
13	Institutio Oratoria	Marcus Fabius Quintilianus	143	1
14	De Officiis	Marcus Tullius Cicero	479	1
15	De Legibus	Marcus Tullius Cicero	241	2
16	De Finibus Bonorum et Malorum	Marcus Tullius Cicero	492	1
17	De Republica	Marcus Tullius Cicero	399	1
18	Noctes Atticae	Aulus Gellius	269	2
19	De Agricultura	Marcus Porcius Cato	236	2
20	De Divinatione	Marcus Tullius Cicero	157	2
1	De Bello Gallico	Gaius Julius Caesar	464	2
2	Metamorphoses	Publius Ovidius Naso	159	1
3	Aeneis	Publius Vergilius Maro	227	1
4	Ab Urbe Condita	Titus Livius	442	2
\.


--
-- TOC entry 3331 (class 0 OID 16587)
-- Dependencies: 217
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, email, password, role) FROM stdin;
1	sebastiano.ferri@gmail.com	$2a$10$heDeRmKnw1pOVutBzCLiMuJ0uFR5ToEWVgMDEwFWIN3xBfyvC02Ou	user
2	giorgiocarloroberto.albieri@gmail.com	$2a$10$7f95R4P2c12MXFSUPaXUyOuKHbf8.pFUiUml3jkDRO.c4a9DYChjS	admin
3	roberto.baggio@hotmail.it	$2a$10$GpgOWLWS8VeHXh0XMzRyhuG3U/kNMZtKSedLkV54KnjQKDUUW5Lrq	user
4	francesco.totti@gmail.it	$2a$10$lK1Ai/IRTIqN0.BbVeifH.c8R7dsaPtY7/hYmqH8tIu3cTZtan2Sm	user
5	alex.delpiero@juve.it	$2a$10$YSJHOk2RgYUydepqRbNah.0AEIO4JsYVn17P8p0/jtqutjtTC274K	user
\.


--
-- TOC entry 3339 (class 0 OID 0)
-- Dependencies: 215
-- Name: books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.books_id_seq', 21, false);


--
-- TOC entry 3340 (class 0 OID 0)
-- Dependencies: 216
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 5, true);


--
-- TOC entry 3181 (class 2606 OID 16408)
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- TOC entry 3183 (class 2606 OID 16596)
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- TOC entry 3185 (class 2606 OID 16594)
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


-- Completed on 2025-09-17 17:41:22

--
-- PostgreSQL database dump complete
--

\unrestrict eVfo24ztCrjq5oR5h4R0FuvgCXc2hThwtUD7gGMwP4d5ltbVgmbhcQoZ9VbVp0D

