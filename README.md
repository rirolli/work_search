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
_Aggiornato il 27/11/2026 20:05 (ora di Roma) · 135 offerte in linea con il profilo su 11706 scaricate · **4 alta**, 43 media, 88 bassa · **28 nuove**_

Filtri e ordinamento: [`docs/index.html`](docs/index.html) · Excel: [`docs/offerte.xlsx`](docs/offerte.xlsx)

| Priorità | Score | | Azienda | Ruolo | Sede | Pubblicata |
|---|---:|---|---|---|---|---|
| Alta | 82 |  | Protolabs | [Senior Machine Learning Engineer - AI Technical Lead](https://www.arbeitnow.nl/jobs/companies/protolabs/senior-machine-learning-engineer-ai-technical-lead-amsterdam-283684) | Amsterdam | 2026-11-25 |
| Alta | 81 |  | ICEYE | [Senior Software Engineer - SAR Processing Platform](https://jobs.ashbyhq.com/iceye/3662acb4-316e-4b9b-aa2d-71ba9f84a85d) | Espoo | 2026-11-13 |
| Alta | 78 |  | ICEYE | [Engineering Team Lead, Flood Observation](https://jobs.ashbyhq.com/iceye/2b6f4ed5-daee-4c4e-aa8f-1a441b600c45) | Espoo | 2026-11-23 |
| Alta | 78 |  | ESA | [Earth Observation Campaign Scientist](https://jobs.esa.int/job/Noordwijk-Earth-Observation-Campaign-Scientist/1458757833/?feedId=null&utm_source=J2WRSS&utm_medium=rss&utm_campaign=J2W_RSS) | Noordwijk, NL | 2026-11-19 |
| Media | 73 |  | ESA | [ERS-EO/Downstream Ground System Engineer](https://jobs.esa.int/job/Frascati-ERS-EODownstream-Ground-System-Engineer/1458807533/?feedId=null&utm_source=J2WRSS&utm_medium=rss&utm_campaign=J2W_RSS) | Frascati, IT | 2026-11-19 |
| Media | 73 |  | ICEYE | [Software Engineer - SAR Processing Platform](https://jobs.ashbyhq.com/iceye/de792d17-acf3-4c8a-bf3f-f582ad84ea0a) | Espoo | 2026-11-13 |
| Media | 73 |  | ICEYE | [Mission Management Team Lead](https://jobs.ashbyhq.com/iceye/bcb9940b-58cd-4c23-8671-7b67f011d07d) | Espoo | 2026-11-11 |
| Media | 73 |  | ICEYE | [AIT Engineering Team Lead](https://jobs.ashbyhq.com/iceye/471cf438-9c29-4167-9c18-be4cc5f4e115) | Lisbon / Warsaw | 2026-08-14 |
| Media | 71 |  | Jobgether | [Senior Machine Learning Engineer, FeedForward 3D/4D Reconstruction](https://jobs.lever.co/jobgether/1b3f4152-df24-4146-a1fb-bcc9ce2c1c7c) | Switzerland / Netherlands / France / Spain / Ireland | 2026-11-25 |
| Media | 71 |  | Open Cosmos | [Edge Computing Software Engineer, Earth Observation](https://careers.open-cosmos.com/jobs/8478374-edge-computing-software-engineer-earth-observation) | Didcot, United Kingdom | 2026-09-30 |
| Media | 70 |  | ICEYE | [Staff Image Processing Engineer](https://jobs.ashbyhq.com/iceye/5613d6f6-24be-4ebf-9b0e-a2284fd21650) | Warsaw | 2026-11-24 |
| Media | 70 |  | Leaf Space | [Project Manager – ESA/EU/ASI Space Projects](https://leafspace.teamtailor.com/jobs/8638002-project-manager-esa-eu-asi-space-projects) | Lomazzo, Italy | 2026-10-29 |
| Media | 68 | 🆕 | Jobgether | [Senior Engineering Team Lead](https://jobs.lever.co/jobgether/1764f255-9d8d-44c6-bcf5-1d666a18d2ca) | Spain / France | 2026-11-26 |
| Media | 68 |  | Safran.AI | [Senior Deep Learning Scientist (H/F)](https://jobs.lever.co/safran-ai/f272d7d4-d63f-4d13-947f-b572fd974d80) | Paris | 2026-11-13 |
| Media | 68 |  | Wayve | [Principal Machine Learning Engineer GAIA](https://jobs.ashbyhq.com/wayve/e781b706-4494-472d-8831-74140ad1a10e) | London, United Kingdom | 2026-10-06 |
| Media | 68 |  | Safran.AI | [Senior Deep Learning Scientist - IMINT (H/F)](https://jobs.lever.co/safran-ai/5875c7a8-54f8-4e6f-bc6e-1e3a5a3e8ae8) | Paris | 2026-09-14 |
| Media | 67 |  | Hydrosat | [Engineering Manager, Applied Analytics](https://careers.hydrosat.com/o/engineering-manager-applied-analytics) | Remote job | 2026-11-11 |
| Media | 67 |  | Wayve | [Senior Machine Learning Engineer, ADAS](https://jobs.ashbyhq.com/wayve/e43c1c44-787e-471b-89fb-ff986c3ab588) | London, United Kingdom | 2026-10-16 |
| Media | 64 |  | LiveEO | [Geospatial Software Engineer (f/m/x)](https://liveeo-gmbh.jobs.personio.com/job/2864893) | LiveEO GmbH Berlin Office (Hybrid) | 2026-10-28 |
| Media | 64 |  | Loft Orbital | [Onboard Software Engineer - Peripheral Services](https://jobs.lever.co/loftorbital/e2428032-ee1a-4e34-a331-914aea52a93f) | Toulouse, Occitanie | 2026-09-22 |
| Media | 64 |  | Loft Orbital | [Onboard Software Engineer - Hyperdrive](https://jobs.lever.co/loftorbital/9182f41f-2d8c-48ff-b914-49a408426e9e) | Toulouse, Occitanie | 2026-09-22 |
| Media | 62 |  | Airbus | [Synthetic Aperture Radar (SAR) Performance Engineer](https://ag.wd3.myworkdayjobs.com/Airbus/job/Getafe-Area/Synthetic-Aperture-Radar--SAR--Performance-Engineer_JR10448575-1) | Getafe Area | 2026-11-04 |
| Media | 61 |  | ICEYE | [Senior MLOps Engineer, ISR](https://jobs.ashbyhq.com/iceye/89154c8e-6501-4012-b71b-8c6711588ffb) | Berlin | 2026-10-05 |
| Media | 60 |  | Airbus | [Earth Observation Satellite System Engineer](https://ag.wd3.myworkdayjobs.com/Airbus/job/Getafe-Area/Earth-Observation-Satellite-System-Engineer_JR10423783) | Getafe Area | Posted 30+ Days Ago |
| Media | 59 |  | Jobgether | [Hands-on Engineering Manager DevOps (GCP/Kubernetes)](https://jobs.lever.co/jobgether/4581e1e7-fe7c-44f6-9f9f-f6fd688c0b91) | Netherlands / UK / Spain / South Africa / Germany / Italy | 2026-11-20 |
| Media | 59 |  | Planet | [Senior Technical Program Manager, Solutions](https://job-boards.greenhouse.io/planetlabs/jobs/8306456) | Berlin, Germany | 2026-10-27 |
| Media | 58 |  | LiveEO | [Senior Project Manager (f/m/x)](https://liveeo-gmbh.jobs.personio.com/job/2873911) | LiveEO GmbH Berlin (Onsite) | 2026-11-03 |
| Media | 57 |  | Jobgether | [Senior Research Scientist 3D Reconstruction & Spatial AI](https://jobs.lever.co/jobgether/f5a0241f-e3f5-41c5-a2d2-853ec0543c58) | Switzerland | 2026-11-25 |
| Media | 57 |  | Helsing | [Staff AI Research Engineer -  3D Computer Vision](https://helsing.ai/jobs/5028718101?gh_jid=5028718101) | Munich - Berlin - London - Paris | 2026-11-06 |
| Media | 57 |  | Wayve | [Senior Machine Learning Engineer, AV Core](https://jobs.ashbyhq.com/wayve/c65879b5-065c-4c29-87ae-3b74f4be2f98) | London, United Kingdom | 2026-10-02 |
| Media | 57 |  | Safran.AI | [Engineering Manager Frontend (H/F)](https://jobs.lever.co/safran-ai/04678422-8092-44ca-ab22-fa82bc09946c) | Paris / Rennes | 2026-09-09 |
| Media | 56 |  | Jobgether | [Engineering Team Lead - AI First](https://jobs.lever.co/jobgether/183023e1-760d-4bcd-a49e-b501078e61d8) | Portugal / Italy / Ireland / Netherlands / Chile / South Africa / Swi… | 2026-11-26 |
| Media | 56 |  | Wayve | [Principal Research Scientist – Multimodal Reasoning & Embodied AI](https://jobs.ashbyhq.com/wayve/a4bd512b-36de-4a28-96b4-80e6b056ee4b) | Sunnyvale, California USA / London, United Kingdom | 2026-11-10 |
| Media | 56 |  | Wayve | [Engineering Manager, Runtime Platform, Robot Software](https://jobs.ashbyhq.com/wayve/34758430-ea52-4ebd-a850-0ea518a41385) | London, United Kingdom | 2026-10-02 |
| Media | 55 |  | Wayve | [Engineering Manager, Experiment Platform, Robot Software](https://jobs.ashbyhq.com/wayve/438bc870-6abc-42a2-be78-1bc746a56a90) | London, United Kingdom | 2026-10-19 |
| Media | 55 |  | ICEYE | [Program Manager - Germany](https://jobs.ashbyhq.com/iceye/4470b416-a308-41a3-a967-293389b62beb) | Berlin / Munich | 2026-10-06 |
| Media | 54 | 🆕 | CesiumAstro | [Senior Machine Learning Engineer](https://www.arbeitnow.com/jobs/companies/cesiumastro/senior-machine-learning-engineer-munich-germany-5311) | Munich, Germany | 2026-11-27 |
| Media | 54 |  | ESA | [FORUM Project Manager](https://jobs.esa.int/job/Noordwijk-FORUM-Project-Manager/1458816433/?feedId=null&utm_source=J2WRSS&utm_medium=rss&utm_campaign=J2W_RSS) | Noordwijk, NL | 2026-11-19 |
| Media | 52 |  | ICEYE | [Lead, AI Hardware Design Skunkworks](https://jobs.ashbyhq.com/iceye/7270373f-9e9c-4608-b492-39afa380fb48) | Valencia | 2026-11-11 |
| Media | 52 |  | Wayve | [Tech Lead - Remote Assist, Robotaxi](https://jobs.ashbyhq.com/wayve/bdf6686e-52e9-45b6-9cf3-3de9645a9b04) | London, United Kingdom | 2026-10-04 |
| Media | 51 | 🆕 | Granola | [Engineering Manager](https://www.arbeitnow.co.uk/jobs/companies/granola/engineering-manager-old-street-329132) | London | 2026-11-27 |
| Media | 51 |  | Jobgether | [Senior Software Engineer / Tech Lead (AI)](https://jobs.lever.co/jobgether/9990af7e-a826-45a6-aa37-a54441e83fe7) | Switzerland | 2026-11-26 |
| Media | 51 |  | Jobgether | [Engineering Manager (Payments)](https://jobs.lever.co/jobgether/09504694-cb9a-4c44-bd6f-8153206320ef) | Portugal / Romania / Poland / Bulgaria | 2026-11-25 |
| Media | 51 | 🆕 | Planet | [Senior Director of Engineering - Civil Government and Commercial](https://job-boards.greenhouse.io/planetlabs/jobs/8324235) | Haarlem, Netherlands / Ljubljana, Slovenia / Berlin, Germany | 2026-10-27 |
| Media | 50 |  | Jobgether | [Senior Data Scientist (Freelance)](https://jobs.lever.co/jobgether/9b0d83cb-b8fc-46ee-b11a-930c7611feb6) | Switzerland | 2026-11-26 |
| Media | 50 |  | Planet | [Technical Program Manager, Solutions](https://job-boards.greenhouse.io/planetlabs/jobs/8396499) | Ljubljana, Slovenia | 2026-11-24 |
| Media | 50 |  | Wayve | [Engineering Manager - Compute and Data Platform](https://jobs.ashbyhq.com/wayve/6ac4d9d0-f0e2-47a1-b0ef-712c9bde1d4d) | London, United Kingdom | 2026-11-04 |
| Bassa | 49 | 🆕 | Jobgether | [Staff DevOps Engineer, AI Security](https://jobs.lever.co/jobgether/1f1d710f-a2b1-4f8e-a997-bfab528dd22e) | Netherlands | 2026-11-27 |
| Bassa | 49 |  | QinetiQ | [Naimuri - Technical Lead Defence](https://careers.qinetiq.com/job/Manchester-Naimuri-Technical-Lead-Defence-Engl/1452198533/?feedId=null&utm_source=J2WRSS&utm_medium=rss&utm_campaign=J2W_RSS) | Manchester, England, United Kingdom | 2026-11-26 |
| Bassa | 49 |  | ICEYE | [AI Solutions Engineer](https://jobs.ashbyhq.com/iceye/9dc49b37-2a22-4e94-9af3-d4312f366cd8) | Espoo | 2026-10-28 |
| Bassa | 49 | 🆕 | ICEYE | [Forward Deployed AI Engineer, Finland](https://jobs.ashbyhq.com/iceye/b6afce4f-317c-4154-82c6-dcd5fa9895ae) | Espoo | 2026-10-19 |
| Bassa | 49 |  | ICEYE | [Forward Deployed AI Engineer, Portugal](https://jobs.ashbyhq.com/iceye/2ef68fcc-6ce8-476e-8c94-6e7b099842e9) | Lisbon / Portugal | 2026-10-15 |
| Bassa | 49 |  | Wayve | [Principal Machine Learning Engineer, Driving Product](https://jobs.ashbyhq.com/wayve/9e02479e-88c0-4e01-bb43-993d9ad5d5c7) | London, United Kingdom | 2026-10-15 |
| Bassa | 49 |  | Planet | [Senior Engineering Program Manager](https://job-boards.greenhouse.io/planetlabs/jobs/8221607) | Berlin, Germany | 2026-09-23 |
| Bassa | 49 |  | ICEYE | [Forward Deployed AI Engineer, Germany](https://jobs.ashbyhq.com/iceye/f1923306-9ffa-4b2f-bee3-c7472c711167) | Munich / Neuss / Berlin | 2026-08-21 |
| Bassa | 48 | 🆕 | Jobgether | [Head of Data & Analytics (iGaming)](https://jobs.lever.co/jobgether/076bcfae-d0ba-4b59-b7d4-3523fc5a00a7) | Belgium / Ireland / France / UK / Netherlands / Switzerland / Romania… | 2026-11-27 |
| Bassa | 48 |  | intrinsicrobotics | [Senior Research Scientist 3D Reconstruction & Spatial AI](https://www.arbeitnow.ch/jobs/companies/intrinsicrobotics/senior-research-scientist-3d-reconstruction-spatial-ai-zurich-143119) | Zurich | 2026-11-24 |
| Bassa | 47 |  | Jobgether | [Media Buying Team Lead](https://jobs.lever.co/jobgether/0426c415-dad1-4a04-bead-94871f112334) | Spain | 2026-11-25 |
| Bassa | 47 |  | Jobgether | [Implementation Team Lead](https://jobs.lever.co/jobgether/2386ecd0-a540-4ae4-85fe-9d3704ae4a2f) | Ireland / UK | 2026-11-25 |
| Bassa | 47 |  | Jobgether | [Engineering Manager, Cloud Platform](https://jobs.lever.co/jobgether/42d18a30-7f04-4ba0-958a-6254afd330ea) | Netherlands | 2026-11-25 |
| Bassa | 47 |  | Safran.AI | [Engineering Manager DevOps (H/F)](https://jobs.lever.co/safran-ai/9304b7d7-69c2-4e54-80cf-8cdd47fdb007) | Paris | 2026-11-24 |
| Bassa | 47 | 🆕 | Jobgether | [Java Tech Lead](https://jobs.lever.co/jobgether/c16948b6-7ad4-46f1-a945-cdac8c01ec51) | Netherlands | 2026-11-21 |
| Bassa | 46 | 🆕 | Jobgether | [Senior Data Scientist - Research EU/UK (m/f/x)](https://jobs.lever.co/jobgether/30efe880-0fd1-49b9-8181-5ac2c4b44185) | Germany | 2026-11-27 |
| Bassa | 45 | 🆕 | Develocity | [Senior Engineering Manager](https://www.arbeitnow.co.uk/jobs/companies/develocity/senior-engineering-manager-united-kingdom-44041) | United Kingdom / Switzerland | 2026-11-27 |
| Bassa | 45 |  | Thales | [Digital Ground Architect for Earth Observation](https://thales.wd3.myworkdayjobs.com/Careers/job/Toulouse/Digital-Ground-Architect-for-Earth-Observation_R0322925-1) | Toulouse | 2026-10-29 |
| Bassa | 45 |  | Wayve | [Machine Learning Engineer, Synthetic Data](https://jobs.ashbyhq.com/wayve/b52585d5-ed8a-4ec8-90f5-2cd06b766e73) | London, United Kingdom | 2026-10-28 |
| Bassa | 45 |  | Wayve | [Tech Lead - Autonomy Performance, Robotaxi](https://jobs.ashbyhq.com/wayve/e457992b-17a7-44eb-8e08-0482a5632386) | London, United Kingdom | 2026-10-09 |
| Bassa | 44 |  | Open Cosmos | [Project Manager, Software Engineering](https://careers.open-cosmos.com/jobs/8751178-project-manager-software-engineering) | Didcot, United Kingdom / Pallini, Greece / Barcelona, Spain / Coimbra… | 2026-11-16 |
| Bassa | 44 |  | Loft Orbital | [Full Stack Engineer - AI Marketplace](https://jobs.lever.co/loftorbital/f5f9edbc-3e8c-4eb3-b757-25cfcbf8e7fa) | Toulouse, Occitanie | 2026-09-07 |
| Bassa | 43 | 🆕 | caronsale | [Engineering Manager – Auctions (m/w/d)](https://www.arbeitnow.com/jobs/companies/caronsale/engineering-manager-auctions-munich-319655) | Munich | 2026-11-27 |
| Bassa | 43 | 🆕 | Deliveroo | [Engineering Manager - Chat Platform](https://www.arbeitnow.co.uk/jobs/companies/deliveroo/engineering-manager-chat-platform-london-312014) | London | 2026-11-27 |
| Bassa | 43 |  | Jobgether | [Team Lead Software Engineering](https://jobs.lever.co/jobgether/8d47bfd3-1c6b-4a6e-bf54-75f75c11911f) | Germany | 2026-11-25 |
| Bassa | 43 |  | Proton | [Senior Engineering Manager](https://www.arbeitnow.ch/jobs/companies/proton/senior-engineering-manager-geneva-paris-france-16972) | Geneva, Paris, France | 2026-11-24 |
| Bassa | 43 |  | Jobgether | [Director of Engineering](https://jobs.lever.co/jobgether/1f07ccda-6d36-464d-91f3-656a88eabb72) | Germany / Switzerland / France / UK / Spain / Ireland / Netherlands | 2026-11-23 |
| Bassa | 43 |  | Wayve | [Lead Technical Program Manager, AI Platform](https://jobs.ashbyhq.com/wayve/002821a8-dd53-4737-be19-dc3270f4d39b) | London, United Kingdom | 2026-10-05 |
| Bassa | 43 |  | Wayve | [Lead Technical Program Manager, Data](https://jobs.ashbyhq.com/wayve/192550ad-e0f6-42d1-a704-414671527a05) | Sunnyvale, California USA / London, United Kingdom | 2026-10-02 |
| Bassa | 43 |  | Astroscale | [Senior Computer Vision Engineer](https://astroscale.bamboohr.com/careers/689) | Oxfordshire UK, Didcot | non indicata |
| Bassa | 42 | 🆕 | Skdse | [Engineering Manager Mission Autonomy (All Genders)](https://www.arbeitnow.com/jobs/companies/skdse/engineering-manager-mission-autonomy-all-genders-berlin-munich-bavaria-363144) | Berlin; Munich, Bavaria | 2026-11-27 |
| Bassa | 42 |  | Channable | [Senior Engineering Manager - Infra Platform Group](https://www.arbeitnow.nl/jobs/companies/channable/senior-engineering-manager-infra-platform-group-utrecht-157004) | Utrecht, Utrecht, Netherlands | 2026-11-25 |
| Bassa | 42 |  | Jobgether | [Senior Regulatory Engineer, AI/ML SaMD](https://jobs.lever.co/jobgether/031238a4-28d6-4c9a-846a-81057065dc1a) | Spain / France / Switzerland | 2026-11-25 |
| Bassa | 42 |  | Airbus | [SAR Systems Engineer](https://ag.wd3.myworkdayjobs.com/Airbus/job/Portsmouth/Systems-Engineer_JR10433681) | Portsmouth | 2026-11-25 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/124c99a8-505d-469c-9362-fd7e4445b670) | Portland | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/dd000db9-8a8b-4403-8354-5bbe5ee1da9c) | Raleigh | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/e62d0897-4c4f-4ccd-9da3-d92490289d4a) | Indianapolis | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/eae75fd8-8d6d-4df1-9ac4-9804945e9a48) | Dallas | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/a812c7ad-af52-4de5-8a44-7670d6a5c1a4) | Phoenix | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/7c7ec6e1-463c-4cf9-aa41-fa259f275682) | Philadelphia | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/a448cf4d-21e5-4c36-b870-50885cc7a959) | Miami | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/c79eadc8-daa6-4c5d-be05-8766c5ebe49e) | San Diego | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/72c5f735-3252-48aa-b27c-a7f960b3339b) | Atlanta | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Principal Data Center & Cloud Engineering Lead and Strategist - Remot…](https://jobs.lever.co/jobgether/22770a24-38e2-433e-a851-20c871ccc812) | Chicago | 2026-11-24 |
| Bassa | 42 |  | Jobgether | [Senior Project Manager](https://jobs.lever.co/jobgether/1a7fedb0-00cb-4157-bb8e-93646a8300ea) | Romania / France / Ireland / Germany / Portugal / Spain / Switzerland… | 2026-11-24 |
| Bassa | 42 |  | Wayve | [Staff Machine Learning Engineer - Ops](https://jobs.ashbyhq.com/wayve/929b2cb4-a3da-4038-ab3c-62cf6f220adb) | London, United Kingdom | 2026-10-02 |
| Bassa | 42 |  | Hydrosat | [AI Native Product Delivery Manager](https://careers.hydrosat.com/o/ai-native-product-delivery-manager) | Luxembourg, Luxembourg, Luxembourg | 2026-08-05 |
| Bassa | 39 | 🆕 | Jobgether | [Senior Applied Scientist - Remote (Portland)](https://jobs.lever.co/jobgether/e7d21bae-6600-4af9-a50e-ca9baa2c507b) | Portland | 2026-11-27 |
| Bassa | 39 | 🆕 | Jobgether | [Senior Applied Scientist - Remote (Philadelphia)](https://jobs.lever.co/jobgether/f3e9b1f4-d23f-45c6-8fcc-0ac1f617cd1c) | Philadelphia | 2026-11-27 |
| Bassa | 39 |  | Jobgether | [Senior AI Agent Engineer – Banking (AI Native)](https://jobs.lever.co/jobgether/6631ceec-b238-49bb-8be1-76456539e8f4) | Spain | 2026-11-25 |
| Bassa | 39 |  | Planet | [Program Manager, German Civil Government](https://job-boards.greenhouse.io/planetlabs/jobs/8226128) | Berlin, Germany | 2026-09-30 |
| Bassa | 38 |  | Airbus | [Senior Artificial Intelligence/Machine Learning Engineer (d/f/m)](https://ag.wd3.myworkdayjobs.com/Airbus/job/Ulm-Area/Senior-Artificial-Intelligence-Machine-Learning-Engineer--d-f-m-_JR10443053-1) | Ulm Area | 2026-11-24 |
| Bassa | 38 |  | QinetiQ | [Technical Lead](https://careers.qinetiq.com/job/Manchester-Technical-Lead-Engl/1266596701/?feedId=null&utm_source=J2WRSS&utm_medium=rss&utm_campaign=J2W_RSS) | Manchester, England, United Kingdom | 2026-11-08 |
| Bassa | 38 |  | Thales | [Artificial Intelligence Space Engineer](https://thales.wd3.myworkdayjobs.com/Careers/job/Roma/Artificial-Intelligence-Space-Engineer_R0210027) | Roma | Posted 30+ Days Ago |
| Bassa | 37 | 🆕 | databricks | [Forward Deployed Software Engineering Manager](https://www.arbeitnow.co.uk/jobs/companies/databricks/forward-deployed-software-engineering-manager-london-66879) | London | 2026-11-27 |
| Bassa | 37 | 🆕 | Jobgether | [Senior IT Project Manager – Laboratory Systems in Pharma](https://jobs.lever.co/jobgether/2269f9ba-a397-4416-9074-ffb8b6463709) | Portugal | 2026-11-26 |
| Bassa | 37 |  | Jobgether | [Senior Project Manager (French Speaker)](https://jobs.lever.co/jobgether/26d91a34-f86b-498f-a367-b2bff5439841) | France | 2026-11-23 |
| Bassa | 37 |  | Wayve | [Staff Software Engineer, AI Libraries](https://jobs.ashbyhq.com/wayve/000a7973-cec6-401c-818a-06781f1f6644) | London, United Kingdom | 2026-10-14 |
| Bassa | 37 |  | Safran.AI | [AI Engineer (H/F)](https://jobs.lever.co/safran-ai/3eedd99a-7f79-4bff-8b8c-46b195ead6c3) | Paris | 2026-10-13 |
| Bassa | 37 |  | Wayve | [Software Engineer, AI Libraries](https://jobs.ashbyhq.com/wayve/82f389e1-b856-4157-ab24-f83f467329e5) | London, United Kingdom | 2026-10-05 |
| Bassa | 37 |  | Wayve | [Staff Validation Systems Engineer, AI Validation](https://jobs.ashbyhq.com/wayve/eef6a32a-f223-4072-8db0-6bf1ca03f3c7) | London, United Kingdom | 2026-10-03 |
| Bassa | 37 |  | Thales | [Ingénieur Docteur en traitement du signal RADAR-SAR H/F](https://thales.wd3.myworkdayjobs.com/Careers/job/Elancourt/Ingnieur-Docteur-en-traitement-du-signal-RADAR-SAR-H-F_R0324957-1) | Elancourt | Posted 30+ Days Ago |
| Bassa | 36 | 🆕 | Fireworks | [Applied Machine Learning Engineer, EMEA](https://www.arbeitnow.co.uk/jobs/companies/fireworks/applied-machine-learning-engineer-emea-london-219855) | London | 2026-11-27 |
| Bassa | 36 | 🆕 | Deliveroo | [Senior Machine Learning Engineer - Recommendations](https://www.arbeitnow.co.uk/jobs/companies/deliveroo/senior-machine-learning-engineer-recommendations-london-267352) | London | 2026-11-27 |
| Bassa | 36 | 🆕 | Deliveroo | [Senior Machine Learning Engineer - Growth & Personalisation](https://www.arbeitnow.co.uk/jobs/companies/deliveroo/senior-machine-learning-engineer-growth-personalisation-london-498446) | London | 2026-11-27 |
| Bassa | 36 |  | Jobgether | [Principal Software Development Engineer in Test (AI-Driven Quality En…](https://jobs.lever.co/jobgether/25c11f2e-10b6-46fe-9c54-90a705c7ce43) | Portugal / Romania / Belgium / Netherlands / Switzerland / UK / Irela… | 2026-11-26 |
| Bassa | 36 | 🆕 | Airbus | [Thèse CIFRE : Réduction de dimension non linéaire et Machine Learning…](https://ag.wd3.myworkdayjobs.com/Airbus/job/Paris-Area/Thse-CIFRE---Rduction-de-dimension-non-linaire-et-Machine-Learning-pour-la-simulation-numrique--H-F-_JR10459513-1) | Paris Area | 2026-11-26 |
| Bassa | 36 |  | Airbus | [AI & Machine Learning Scientist – Materials Research (d/m/w)](https://ag.wd3.myworkdayjobs.com/Airbus/job/Mnchen-Area/AI---Machine-Learning-Scientist---Materials-Research--d-m-w-_JR10458316) | München Area | 2026-11-25 |
| Bassa | 36 |  | QinetiQ | [Naimuri - Senior Data Scientist](https://careers.qinetiq.com/job/Manchester-Naimuri-Senior-Data-Scientist-Engl/1280625801/?feedId=null&utm_source=J2WRSS&utm_medium=rss&utm_campaign=J2W_RSS) | Manchester, England, United Kingdom | 2026-11-17 |
| Bassa | 35 | 🆕 | Jobgether | [Tech Lead IA & Arquitecto/a software](https://jobs.lever.co/jobgether/7916b7a6-7da0-4f45-bba6-dff1d67922e3) | Spain | 2026-11-26 |
| Bassa | 35 |  | Jobgether | [Senior Engineer, Data and ML Platform](https://jobs.lever.co/jobgether/d5b7d06c-57c5-4f1b-94ec-290a5c5ed0b6) | Switzerland | 2026-11-26 |
| Bassa | 35 |  | Jobgether | [Senior Artificial Intelligence Specialist](https://jobs.lever.co/jobgether/050b324e-7524-4fc8-9bd7-9a40da81e38b) | Italy / Ireland / Switzerland / Spain / Portugal / Netherlands / Roma… | 2026-11-23 |
| Bassa | 34 |  | Jobgether | [Principal Applied AI Engineer / Forward-Deployed Engineer - Remote (C…](https://jobs.lever.co/jobgether/bd8c6491-ea38-412a-b309-6801b1b25e3e) | Chicago | 2026-11-25 |
| Bassa | 34 |  | Airbus | [Geospatial Commercial & Bid Manager](https://ag.wd3.myworkdayjobs.com/Airbus/job/Barcelona/Geospatial-Commercial---Bid-Manager_JR10350385) | Barcelona | Posted 30+ Days Ago |
| Bassa | 33 | 🆕 | Jobgether | [Program Director/ Project Manager - Remote (Chicago)](https://jobs.lever.co/jobgether/ae66d2b0-a01a-4c4d-b2d3-052614f501bd) | Chicago | 2026-11-26 |
| Bassa | 32 | 🆕 | Safran.AI | [Junior Deep Learning Scientist (H/F)](https://jobs.lever.co/safran-ai/934211e0-993c-49ec-998b-51fb586e32f8) | Paris | 2026-11-27 |
| Bassa | 32 | 🆕 | 7NXT GmbH | [Senior Data Engineer- Data & AI Platform (f/m/d)](https://www.arbeitnow.com/jobs/companies/7nxt-gmbh/senior-data-engineer-data-ai-platform-berlin-127467) | Berlin | 2026-11-27 |
| Bassa | 32 |  | Jobgether | [Chief Technology Officer (Java / Agentic AI)](https://jobs.lever.co/jobgether/4daf458d-7395-4781-8b02-50b5e5fd34c9) | Italy | 2026-11-26 |
| Bassa | 32 |  | Thales | [Tech Lead IA F/H](https://thales.wd3.myworkdayjobs.com/Careers/job/Vlizy-Villacoublay/Tech-Lead-IA-F-H_R0344884-1) | Vélizy-Villacoublay | 2026-11-18 |
| Bassa | 32 |  | Wayve | [Release Manager - AI Models](https://jobs.ashbyhq.com/wayve/c8056c08-472a-4d64-8d69-a6520083bf80) | London, United Kingdom | 2026-10-21 |
| Bassa | 32 |  | RTX | [Principal Data Scientist & Machine Learning Researcher](https://globalhr.wd5.myworkdayjobs.com/REC_RTX_Ext_Gateway/job/Gloucester-South-Gloucestershire/Principal-Data-Scientist---Machine-Learning-Researcher_01804570) | 3 Locations | Posted 30+ Days Ago |
| Bassa | 30 | 🆕 | Outsight | [Senior Software Engineer - Perception](https://www.arbeitnow.fr/jobs/companies/outsight/senior-software-engineer-perception-paris-465218) | Paris | 2026-11-27 |
| Bassa | 30 | 🆕 | Jobgether | [AI Specialist](https://jobs.lever.co/jobgether/5107c0ec-75cc-4fb2-9f92-a274dfd0e8cc) | Switzerland | 2026-11-26 |
| Bassa | 30 |  | Jobgether | [AI Solution Architect](https://jobs.lever.co/jobgether/21b858fa-ad0e-4574-a2a1-1fa68158b510) | Argentina | 2026-11-25 |
| Bassa | 30 |  | Jobgether | [Forward Deployed Engineer (AI & Digital Experience Platforms)](https://jobs.lever.co/jobgether/0c4b28f6-a430-488a-9dc0-2f5d6dbce8bd) | UK / France / Netherlands / Switzerland / Ireland / Germany / Spain | 2026-11-23 |
| Bassa | 30 | 🆕 | Jobgether | [Senior AI Engineer (Internal Tools & Agents)](https://jobs.lever.co/jobgether/7026821f-1d4a-4b0f-ac75-b885ac6a20d7) | Italy | 2026-11-20 |
| Bassa | 30 |  | Wayve | [Applied Scientist/Machine Learning Engineer, Gaia](https://jobs.ashbyhq.com/wayve/5d3baafc-4af6-47db-91a0-ee5cb2cf3e04) | London, United Kingdom | 2026-10-02 |
| Bassa | 30 |  | ClearSpace | [On-board Data Handling Engineer](https://clearspace.bamboohr.com/careers/81) | London, London, City of | non indicata |
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
