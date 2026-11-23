# work_search – ricerca lavoro automatica (Earth Observation / ML)

Una GitHub Action esegue la ricerca **alle 8:00 e alle 20:00 (ora di Roma)**. A ogni esecuzione:

1. scarica le offerte dalle pagine carriere delle aziende in `config/companies.json` e dagli aggregatori in `config/searches.json`;
2. tiene solo quelle in linea con il profilo (`config/profile.json`) e assegna a ciascuna un punteggio e una priorità;
3. confronta il risultato con le esecuzioni precedenti (`data/state.json`): segna le offerte **nuove** e **rimuove** quelle che le fonti non elencano più;
4. pubblica in `docs/`:
   - `index.html`: pagina autonoma con filtri (testo, priorità, aerospaziale sì/no, solo nuove, solo Italia/Roma) e ordinamento per colonna;
   - `offerte.xlsx`: fogli *Offerte*, *Fonti* (esito di ogni fonte) e *Info* (criteri);
   - `offerte.json`: gli stessi dati in formato JSON.

I file vengono committati nel repository e allegati all'esecuzione come artefatto (`ricerca-lavoro`). Tutte le offerte compaiono anche qui sotto, nella sezione **Risultati**, che la routine riscrive a ogni esecuzione (non modificare il testo tra i due marcatori).

## Risultati

<!-- RISULTATI:INIZIO -->
_Aggiornato il 23/11/2026 12:58 (ora di Roma) · 68 offerte in linea con il profilo su 6167 scaricate · **4 alta**, 21 media, 43 bassa · **68 nuove**_

Filtri e ordinamento: [`docs/index.html`](docs/index.html) · Excel: [`docs/offerte.xlsx`](docs/offerte.xlsx)

| Priorità | Score | | Azienda | Ruolo | Sede | Pubblicata |
|---|---:|---|---|---|---|---|
| Alta | 93 | 🆕 | ICEYE | [Engineering Team Lead (Geospatial Analytics)](https://jobs.ashbyhq.com/iceye/e5082d17-b132-4106-b70d-2a08cdeeb405) | Espoo | 2026-10-30 |
| Alta | 81 | 🆕 | ICEYE | [Senior Software Engineer - SAR Processing Platform](https://jobs.ashbyhq.com/iceye/3662acb4-316e-4b9b-aa2d-71ba9f84a85d) | Espoo | 2026-11-13 |
| Alta | 78 | 🆕 | ICEYE | [Engineering Team Lead, Flood Observation](https://jobs.ashbyhq.com/iceye/2b6f4ed5-daee-4c4e-aa8f-1a441b600c45) | Espoo | 2026-11-23 |
| Alta | 77 | 🆕 | ICEYE | [Senior Geospatial Software Engineer](https://jobs.ashbyhq.com/iceye/4b5957fb-dd07-424f-b1b0-59bbd2864a0f) | Espoo | 2026-11-09 |
| Media | 73 | 🆕 | ICEYE | [Software Engineer - SAR Processing Platform](https://jobs.ashbyhq.com/iceye/de792d17-acf3-4c8a-bf3f-f582ad84ea0a) | Espoo | 2026-11-13 |
| Media | 73 | 🆕 | ICEYE | [Mission Management Team Lead](https://jobs.ashbyhq.com/iceye/bcb9940b-58cd-4c23-8671-7b67f011d07d) | Espoo | 2026-11-11 |
| Media | 73 | 🆕 | ICEYE | [AIT Engineering Team Lead](https://jobs.ashbyhq.com/iceye/471cf438-9c29-4167-9c18-be4cc5f4e115) | Lisbon / Warsaw | 2026-08-14 |
| Media | 64 | 🆕 | Jobgether | [Hands-on Engineering Manager DevOps (GCP/Kubernetes)](https://jobs.lever.co/jobgether/4581e1e7-fe7c-44f6-9f9f-f6fd688c0b91) | Netherlands / UK / Romania / Switzerland / France / Spain / Portugal … | 2026-11-20 |
| Media | 64 | 🆕 | Mytomorrows | [Senior Engineering Manager](https://www.arbeitnow.nl/jobs/companies/mytomorrows/senior-engineering-manager-amsterdam-485449) | Amsterdam HQ | 2026-11-19 |
| Media | 64 | 🆕 | LiveEO | [Geospatial Software Engineer (f/m/x)](https://liveeo-gmbh.jobs.personio.com/job/2864893) | LiveEO GmbH Berlin Office (Hybrid) | 2026-10-28 |
| Media | 64 | 🆕 | Loft Orbital | [Onboard Software Engineer - Peripheral Services](https://jobs.lever.co/loftorbital/e2428032-ee1a-4e34-a331-914aea52a93f) | Toulouse, Occitanie | 2026-09-22 |
| Media | 64 | 🆕 | Loft Orbital | [Onboard Software Engineer - Hyperdrive](https://jobs.lever.co/loftorbital/9182f41f-2d8c-48ff-b914-49a408426e9e) | Toulouse, Occitanie | 2026-09-22 |
| Media | 62 | 🆕 | Airbus | [Synthetic Aperture Radar (SAR) Performance Engineer](https://ag.wd3.myworkdayjobs.com/Airbus/job/Getafe-Area/Synthetic-Aperture-Radar--SAR--Performance-Engineer_JR10448575-1) | Getafe Area | 2026-11-04 |
| Media | 62 | 🆕 | Black Shark | [Software Engineering Manager - Data Plane](https://jobs.eu.lever.co/blackshark/023ce5b4-3800-4b87-ac17-bdc30a6ddbd1) | Graz | 2026-09-24 |
| Media | 61 | 🆕 | ICEYE | [Senior MLOps Engineer, ISR](https://jobs.ashbyhq.com/iceye/89154c8e-6501-4012-b71b-8c6711588ffb) | Berlin | 2026-10-05 |
| Media | 61 | 🆕 | Black Shark | [Software Engineering Manager - Control Plane](https://jobs.eu.lever.co/blackshark/7c292d32-b2d9-429e-94b5-a7907920429d) | Graz | 2026-09-24 |
| Media | 60 | 🆕 | Mytomorrows | [Sr. Engineering Manager](https://www.arbeitnow.nl/jobs/companies/mytomorrows/sr-engineering-manager-amsterdam-167050) | Amsterdam HQ | 2026-11-21 |
| Media | 60 | 🆕 | Airbus | [Earth Observation Satellite System Engineer](https://ag.wd3.myworkdayjobs.com/Airbus/job/Getafe-Area/Earth-Observation-Satellite-System-Engineer_JR10423783) | Getafe Area | Posted 30+ Days Ago |
| Media | 59 | 🆕 | Planet | [Senior Technical Program Manager, Solutions](https://job-boards.greenhouse.io/planetlabs/jobs/8306456) | Berlin, Germany | 2026-10-27 |
| Media | 58 | 🆕 | LiveEO | [Senior Project Manager (f/m/x)](https://liveeo-gmbh.jobs.personio.com/job/2873911) | LiveEO GmbH Berlin (Onsite) | 2026-11-03 |
| Media | 55 | 🆕 | Jobgether | [Machine Learning Engineer](https://jobs.lever.co/jobgether/15d4c9c5-9434-4f11-b44d-182b593664c1) | Poland / Netherlands / Sweden / UK / Romania / Austria / Belgium / Po… | 2026-11-17 |
| Media | 55 | 🆕 | ICEYE | [Program Manager - Germany](https://jobs.ashbyhq.com/iceye/4470b416-a308-41a3-a967-293389b62beb) | Berlin / Munich | 2026-10-06 |
| Media | 52 | 🆕 | Laelaps | [Research Scientist, Foundation Model (World Model & Site Memory)](https://www.arbeitnow.ch/jobs/companies/laelaps/research-scientist-foundation-model-world-model-site-memory-zurich-448097) | Zürich, Switzerland | 2026-11-22 |
| Media | 52 | 🆕 | ICEYE | [Lead, AI Hardware Design Skunkworks](https://jobs.ashbyhq.com/iceye/7270373f-9e9c-4608-b492-39afa380fb48) | Valencia | 2026-11-11 |
| Media | 51 | 🆕 | Planet | [Senior Director of Engineering - Civil Government and Commercial](https://job-boards.greenhouse.io/planetlabs/jobs/8324235) | Haarlem, Netherlands / Ljubljana, Slovenia / Berlin, Germany | 2026-10-27 |
| Bassa | 49 | 🆕 | Vestiairecollective | [Senior Data Scientist - Applied Scientist](https://www.arbeitnow.com/jobs/companies/vestiairecollective/senior-data-scientist-applied-scientist-berlin-487392) | Berlin | 2026-11-23 |
| Bassa | 49 | 🆕 | mirakllabs | [MLOps Engineer (CDI)](https://www.arbeitnow.fr/jobs/companies/mirakllabs/mlops-engineer-cdi-paris-386846) | Paris / Bordeaux | 2026-11-21 |
| Bassa | 49 | 🆕 | ICEYE | [AI Solutions Engineer](https://jobs.ashbyhq.com/iceye/9dc49b37-2a22-4e94-9af3-d4312f366cd8) | Espoo | 2026-10-28 |
| Bassa | 49 | 🆕 | ICEYE | [Forward Deployed AI Engineer, Portugal](https://jobs.ashbyhq.com/iceye/2ef68fcc-6ce8-476e-8c94-6e7b099842e9) | Lisbon / Portugal | 2026-10-15 |
| Bassa | 49 | 🆕 | Planet | [Senior Engineering Program Manager](https://job-boards.greenhouse.io/planetlabs/jobs/8221607) | Berlin, Germany | 2026-09-23 |
| Bassa | 49 | 🆕 | ICEYE | [Forward Deployed AI Engineer, Germany](https://jobs.ashbyhq.com/iceye/f1923306-9ffa-4b2f-bee3-c7472c711167) | Munich / Neuss / Berlin | 2026-08-21 |
| Bassa | 48 | 🆕 | DiliTrust | [Engineering Manager - AI platform](https://www.arbeitnow.fr/jobs/companies/dilitrust/engineering-manager-ai-platform-paris-la-defense-312019) | Paris - La Défense | 2026-11-23 |
| Bassa | 48 | 🆕 | FINN | [Engineering Manager - Cloud & Security (m/f/x)](https://www.arbeitnow.com/jobs/companies/finn/engineering-manager-cloud-security-munich-236210) | Munich | 2026-11-23 |
| Bassa | 47 | 🆕 | Jobgether | [Java Tech Lead](https://jobs.lever.co/jobgether/007bcb66-ef92-4294-be42-f9f9ca6c5b1b) | UK / Germany / France / Switzerland / Ireland / Netherlands / Spain | 2026-11-21 |
| Bassa | 47 | 🆕 | Jobgether | [Senior ML Engineer](https://jobs.lever.co/jobgether/10fd133e-d0a2-43db-bb9d-9f61bfa27550) | Germany / Spain / UK / Netherlands / Italy / France / Ireland / Switz… | 2026-11-16 |
| Bassa | 45 | 🆕 | Thales | [Digital Ground Architect for Earth Observation](https://thales.wd3.myworkdayjobs.com/Careers/job/Toulouse/Digital-Ground-Architect-for-Earth-Observation_R0322925-1) | Toulouse | 2026-10-29 |
| Bassa | 44 | 🆕 | Loft Orbital | [Full Stack Engineer - AI Marketplace](https://jobs.lever.co/loftorbital/f5f9edbc-3e8c-4eb3-b757-25cfcbf8e7fa) | Toulouse, Occitanie | 2026-09-07 |
| Bassa | 43 | 🆕 | Jobgether | [Director of Engineering](https://jobs.lever.co/jobgether/1f07ccda-6d36-464d-91f3-656a88eabb72) | Germany / Switzerland / France / UK / Spain / Ireland / Netherlands | 2026-11-23 |
| Bassa | 42 | 🆕 | Makersite GmbH | [Senior Data Scientist - Research EU/UK (m/f/x)](https://www.arbeitnow.com/jobs/companies/makersite-gmbh/remote-senior-data-scientist-research-eu-uk-berlin-88725) | Remote job | 2026-11-22 |
| Bassa | 41 | 🆕 | Vinteden | [Data Scientist, Trust & Safety](https://www.arbeitnow.nl/jobs/companies/vinteden/data-scientist-trust-safety-amsterdam-87460) | Amsterdam | 2026-11-19 |
| Bassa | 39 | 🆕 | Bridgemaker GmbH | [Tech Lead Marketplace Experience (m/w/d)](https://www.arbeitnow.com/jobs/companies/bridgemaker-gmbh/tech-lead-marketplace-experience-berlin-278949) | Berlin | 2026-11-23 |
| Bassa | 39 | 🆕 | Bridgemaker GmbH | [Tech Lead Grower Experience (m/w/d)](https://www.arbeitnow.com/jobs/companies/bridgemaker-gmbh/tech-lead-grower-experience-berlin-336126) | Berlin | 2026-11-23 |
| Bassa | 39 | 🆕 | Jobgether | [ML Engineer](https://jobs.lever.co/jobgether/0cecd754-7572-48d0-b427-b372c02225c6) | Greece / UK / Spain / France / Belgium / Kenya / Portugal / Luxembour… | 2026-11-18 |
| Bassa | 39 | 🆕 | Planet | [Program Manager, German Civil Government](https://job-boards.greenhouse.io/planetlabs/jobs/8226128) | Berlin, Germany | 2026-09-30 |
| Bassa | 38 | 🆕 | Airbus | [Senior Machine Learning Engineer (d/m/w)](https://ag.wd3.myworkdayjobs.com/Airbus/job/Stuttgart-Area/Senior-Machine-Learning-Engineer--d-m-w-_JR10459155) | Stuttgart Area | 2026-11-18 |
| Bassa | 38 | 🆕 | Thales | [Artificial Intelligence Space Engineer](https://thales.wd3.myworkdayjobs.com/Careers/job/Roma/Artificial-Intelligence-Space-Engineer_R0210027) | Roma | Posted 30+ Days Ago |
| Bassa | 37 | 🆕 | Thales | [Tech Lead IA F/H](https://thales.wd3.myworkdayjobs.com/Careers/job/Vlizy-Villacoublay/Tech-Lead-IA-F-H_R0344884-1) | Vélizy-Villacoublay | 2026-11-18 |
| Bassa | 37 | 🆕 | Thales | [Ingénieur Docteur en traitement du signal RADAR-SAR H/F](https://thales.wd3.myworkdayjobs.com/Careers/job/Elancourt/Ingnieur-Docteur-en-traitement-du-signal-RADAR-SAR-H-F_R0324957-1) | Elancourt | Posted 30+ Days Ago |
| Bassa | 36 | 🆕 | Theodo | [Forward Deployed Engineer - CDI Paris - Theodo AI-Modernisation](https://www.arbeitnow.fr/jobs/companies/theodo/forward-deployed-engineer-cdi-paris-theodo-ai-modernisation-389755) | Paris | 2026-11-23 |
| Bassa | 35 | 🆕 | Jobgether | [Senior Artificial Intelligence Specialist](https://jobs.lever.co/jobgether/050b324e-7524-4fc8-9bd7-9a40da81e38b) | Italy / Ireland / Switzerland / Spain / Portugal / Netherlands / Roma… | 2026-11-23 |
| Bassa | 35 | 🆕 | adjoe | [Senior/Staff ML Scientist - Bidding & Optimization (f/m/d)](https://www.arbeitnow.com/jobs/companies/adjoe/senior-staff-ml-scientist-bidding-optimization-berlin-248715) | Berlin | 2026-11-22 |
| Bassa | 35 | 🆕 | Jobgether | [Senior AI Engineer (Internal Tools & Agents)](https://jobs.lever.co/jobgether/0b7492c8-2270-4764-9185-da7c1176ea81) | Chile / Belgium / South Africa / Croatia / Spain / Germany / Argentin… | 2026-11-20 |
| Bassa | 35 | 🆕 | Flexport | [Engineering Manager, Customs](https://www.arbeitnow.nl/jobs/companies/flexport/engineering-manager-customs-amsterdam-328631) | Amsterdam | 2026-11-18 |
| Bassa | 34 | 🆕 | Airbus | [Geospatial Commercial & Bid Manager](https://ag.wd3.myworkdayjobs.com/Airbus/job/Barcelona/Geospatial-Commercial---Bid-Manager_JR10350385) | Barcelona | Posted 30+ Days Ago |
| Bassa | 33 | 🆕 | Bjak | [Engineering Manager - AI-Native Email App](https://www.arbeitnow.co.uk/jobs/companies/bjak/engineering-manager-ai-native-email-app-london-444103) | United Kingdom | 2026-11-23 |
| Bassa | 33 | 🆕 | Bjak | [Engineering Manager - Stock Trading App](https://www.arbeitnow.co.uk/jobs/companies/bjak/engineering-manager-stock-trading-app-london-6781) | United Kingdom | 2026-11-23 |
| Bassa | 33 | 🆕 | Bjak | [Engineering Manager - Wealth Management App](https://www.arbeitnow.co.uk/jobs/companies/bjak/engineering-manager-wealth-management-app-london-99649) | United Kingdom | 2026-11-23 |
| Bassa | 33 | 🆕 | Bjak | [Engineering Manager - Travel FX App](https://www.arbeitnow.co.uk/jobs/companies/bjak/engineering-manager-travel-fx-app-london-213505) | United Kingdom | 2026-11-23 |
| Bassa | 33 | 🆕 | Bjak | [Engineering Manager - AI Payments App](https://www.arbeitnow.co.uk/jobs/companies/bjak/engineering-manager-ai-payments-app-london-489211) | United Kingdom | 2026-11-23 |
| Bassa | 33 | 🆕 | Bjak | [Engineering Manager](https://www.arbeitnow.co.uk/jobs/companies/bjak/engineering-manager-london-377260) | United Kingdom | 2026-11-23 |
| Bassa | 33 | 🆕 | ILLUIN Technology | [Senior Applied AI Manager (F/H)](https://www.arbeitnow.fr/jobs/companies/illuin-technology/senior-applied-ai-manager-f-h-puteaux-474357) | Paris La Défense | 2026-11-23 |
| Bassa | 33 | 🆕 | Airbus | [Machine Learning Engineer (MLOps & AI Integration)](https://ag.wd3.myworkdayjobs.com/Airbus/job/Albacete/Machine-Learning-Engineer--MLOps---AI-Integration-_JR10454310) | Albacete | 2026-11-05 |
| Bassa | 32 | 🆕 | Jobgether | [Senior Data Scientist – Applied ML & GenAI](https://jobs.lever.co/jobgether/ee0c58b3-e91b-451f-ae67-6e734c6ecf0f) | Italy | 2026-11-18 |
| Bassa | 31 | 🆕 | databricks | [AI Forward Deployed Engineer - Netherlands](https://www.arbeitnow.nl/jobs/companies/databricks/ai-forward-deployed-engineer-netherlands-amsterdam-113124) | Amsterdam | 2026-11-20 |
| Bassa | 30 | 🆕 | Jobgether | [Forward Deployed Engineer (AI & Digital Experience Platforms)](https://jobs.lever.co/jobgether/0c4b28f6-a430-488a-9dc0-2f5d6dbce8bd) | UK / France / Netherlands / Switzerland / Ireland / Germany / Spain | 2026-11-23 |
| Bassa | 30 | 🆕 | mirakllabs | [Senior AI Engineer - Product](https://www.arbeitnow.fr/jobs/companies/mirakllabs/senior-ai-engineer-product-paris-100724) | Paris | 2026-11-21 |
| Bassa | 30 | 🆕 | mirakllabs | [Senior AI Engineer Freelance - Product](https://www.arbeitnow.fr/jobs/companies/mirakllabs/senior-ai-engineer-freelance-product-paris-35952) | Paris | 2026-11-21 |
| Bassa | 30 | 🆕 | Jobgether | [Lead Engineer - AI Platform](https://jobs.lever.co/jobgether/0d3e06d5-2b41-451e-a5f5-293dbababd00) | Ireland / Germany / Spain / Switzerland / UK / Netherlands / France | 2026-11-17 |
<!-- RISULTATI:FINE -->

## Primo avvio

1. **Actions → Ricerca lavoro → Run workflow** per lanciarla subito a mano.
2. Apri il foglio **Fonti** (o la sezione *Stato delle fonti* dell'HTML): le fonti segnate *(endpoint da verificare)* sono identificativi presunti. Se danno errore, correggi `id`/`domain` in `companies.json` o rimuovi la fonte.
3. *(consigliato)* Aggiungi le chiavi gratuite degli aggregatori in **Settings → Secrets and variables → Actions**. Senza chiavi quelle ricerche vengono saltate. Sono soprattutto loro a coprire le aziende italiane senza API pubblica (Telespazio, Thales Alenia Space, ESA, Eni, Enel, Sky…):
   - `ADZUNA_APP_ID`, `ADZUNA_APP_KEY` – registrazione su developer.adzuna.com (copre Italia, Germania, Francia, Olanda, Spagna, UK…);
   - `JOOBLE_API_KEY` – richiesta su jooble.org/api/about.
4. *(facoltativo)* Per vedere l'HTML come sito: **Settings → Pages → Deploy from a branch → `main` / `docs`**. Attenzione: con il repository pubblico la pagina sarà pubblica.

## Fonti supportate

| Tipo | Campi in `sources` | Note |
|---|---|---|
| `greenhouse` | `id` (board token), `region: "eu"` facoltativo | es. Planet `planetlabs`, DeepMind `deepmind`, Scandit `scandit` |
| `lever` | `id`, `region: "eu"` facoltativo | es. Loft Orbital, Black Shark (EU) |
| `ashby` | `id` | es. ICEYE, Overstory |
| `personio` | `id`, `domain: "de"/"com"` | feed XML pubblico, es. LiveEO |
| `recruitee` | `id` | |
| `workable` | `id` (account) | |
| `smartrecruiters` | `id` (company identifier) | |
| `workday` | `domain`, `tenant`, `site`, `queries` | pagine `*.myworkdayjobs.com` (Airbus, Thales); pubblicazione solo relativa |
| `teamtailor` | `domain` | feed `/jobs.rss`, es. Open Cosmos, Kuva Space |
| `pinpoint` | `domain` | es. osapiens |
| `successfactors` | `domain`, `queries` facoltativo | feed RSS di SAP SuccessFactors, es. ESA (`jobs.esa.int`) |
| `bamboohr`, `breezy` | `id` | |
| aggregatori | in `searches.json`: `adzuna`, `jooble`, `arbeitnow`, `remotive` | Adzuna/Jooble con chiave |

Un'azienda senza `sources` viene comunque riconosciuta nei risultati degli aggregatori tramite `name`/`aliases` e riceve settore e rating. `exclude: true` la scarta (è impostato per e-GEOS, il datore di lavoro attuale).

## Aggiungere aziende: scoperta automatica delle fonti

1. Aggiungi l'azienda in `config/candidates.json` (nome, eventuali `slugs` alternativi e `domains` del sito carriere).
2. Lancia **Actions → Scoperta fonti → Run workflow** (campo facoltativo per filtrare per nome).
3. Nel riepilogo dell'esecuzione trovi, per ogni azienda, le API che hanno risposto, con il numero di offerte, alcuni titoli di esempio e lo snippet da copiare in `sources` di `config/companies.json`. Controlla i titoli di esempio: un nome generico può corrispondere a un'altra azienda omonima.

Le aziende senza API pubblica (Telespazio, Leonardo, Thales Alenia Space a parte Workday, Eni, Enel, ECMWF, EUMETSAT, DeepMind…) restano in `companies.json` senza `sources`: vengono riconosciute nei risultati di Adzuna e Jooble, che però richiedono le chiavi.

## Punteggio

- Il titolo deve contenere un termine "core" di `title_required_any` (ML, computer vision, EO, SAR, geospatial…) e nessuno di `title_exclude`.
- In alternativa può essere un ruolo gestionale (`title_manager_any`: team lead, engineering/technical manager, project manager…), ma solo se titolo o descrizione citano un ambito di `manager_requires_any` (ML, AI, EO, remote sensing…).
- Ogni parola chiave vale `peso × 2` nel titolo e `peso` nella descrizione; i punti dalla descrizione hanno un tetto (`description_max`, 30), così le aziende spaziali non salgono solo per il gergo del settore.
- Si aggiungono il peso della sede (Roma 30, Frascati 26, Italia 15, Europa 4–8, remoto 10), il rating dell'azienda e un bonus se l'offerta è recente.
- Le sedi in `location_exclude` (USA e sigle degli stati, India, Giappone…) vengono scartate, a meno che la stessa offerta indichi anche una sede europea.
- Lo stesso annuncio pubblicato in più sedi diventa una sola riga, con le sedi unite.
- Priorità **Alta** ≥ 75, **Media** ≥ 50; sotto 30 l'offerta non compare. Tutte le soglie si cambiano in `profile.json`.

I termini si confrontano come parole intere: `ai` non corrisponde dentro "maintenance".

## Limiti

- **Offerte rimosse**: un'offerta che la fonte non elenca più sparisce da README, HTML ed Excel alla prima esecuzione successiva. Se invece la fonte è in errore, le sue offerte non compaiono in quell'esecuzione ma restano in memoria per 7 giorni (`state_retention_days`), così al ritorno non vengono segnate come nuove.
- **Data di scadenza**: quasi nessuna API la espone, quindi la colonna riporta "non indicata". *Prima vista* indica quando la routine ha trovato l'offerta.
- **LinkedIn** non viene interrogato: non offre un'API pubblica per le offerte.
- Il cron di GitHub è in UTC e può partire con qualche minuto di ritardo. Il workflow pianifica entrambi gli orari (ora solare e legale) e tiene solo quello che corrisponde alle 8 o alle 20 a Roma.
- GitHub sospende i workflow pianificati dopo 60 giorni senza attività nel repository. I commit giornalieri dei risultati di norma bastano; se succede, riattivalo dalla scheda Actions.

## Esecuzione locale

```bash
go test ./...
go run ./cmd/jobsearch -config config -out docs -state data/state.json
go run ./cmd/jobsearch -only iceye      # solo le fonti il cui nome contiene "iceye" (debug)
```

Il programma usa solo la libreria standard di Go (1.22+), senza dipendenze esterne.

## Struttura

```
cmd/jobsearch/        comando principale
cmd/discover/         scoperta delle API delle pagine carriere
internal/sources/     connettori ATS e aggregatori
internal/score/       filtro e punteggio
internal/state/       memoria delle offerte già viste
internal/report/      HTML (template.html) e Excel
internal/xlsx/        writer .xlsx minimale
config/               aziende, ricerche, profilo
data/state.json       stato (aggiornato dalla routine)
docs/                 output (aggiornato dalla routine)
.github/workflows/    job-search.yml (8:00/20:00), discover.yml (manuale)
```
