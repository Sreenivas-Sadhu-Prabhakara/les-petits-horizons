"""Generate a synthetic French-learning corpus as JSONL (one document per line).

Phase-1 placeholder material so the RAG pipeline works end-to-end before real
curriculum lands in Phase 2. Graded notes, vocabulary and example dialogues
across CEFR levels A1–C2. Re-runnable and deterministic.
Output: data/seed/corpus.jsonl  (all docs: audience="french", language="fr")
"""
import json, os, hashlib

OUT_DIR = os.path.join(os.path.dirname(__file__), "seed")
OUT = os.path.join(OUT_DIR, "corpus.jsonl")

DOCS = [
    # --- A1 : salutations, présentations, base ---
    {"title": "Salutations et politesse (A1)", "source_type": "lesson",
     "content": "Pour dire bonjour le matin et la journée : « Bonjour ». Le soir : « Bonsoir ». "
                "De façon familière entre amis : « Salut ». Pour dire au revoir : « Au revoir », "
                "« À bientôt », « À demain ». Les mots de politesse essentiels : « s'il vous plaît » "
                "(formel), « s'il te plaît » (familier), « merci », « de rien », « excusez-moi », « pardon »."},
    {"title": "Se présenter (A1)", "source_type": "lesson",
     "content": "Pour se présenter : « Je m'appelle Marie. » ou « Mon nom est Marie. » Pour demander : "
                "« Comment tu t'appelles ? » (familier) ou « Comment vous appelez-vous ? » (formel). "
                "Dire son origine : « Je viens de France. », « Je suis français / française. » "
                "Demander comment ça va : « Comment ça va ? » — réponses : « Ça va bien, merci. », « Ça va. », « Pas mal. »"},
    {"title": "Les nombres de 0 à 20 (A1)", "source_type": "lesson",
     "content": "0 zéro, 1 un, 2 deux, 3 trois, 4 quatre, 5 cinq, 6 six, 7 sept, 8 huit, 9 neuf, 10 dix, "
                "11 onze, 12 douze, 13 treize, 14 quatorze, 15 quinze, 16 seize, 17 dix-sept, 18 dix-huit, "
                "19 dix-neuf, 20 vingt. Astuce : à partir de 17, on combine dix + le chiffre (dix-sept = 10 + 7)."},
    {"title": "Articles définis et indéfinis (A1)", "source_type": "grammar",
     "content": "Les articles indéfinis : « un » (masculin : un livre), « une » (féminin : une table), "
                "« des » (pluriel : des livres). Les articles définis : « le » (le livre), « la » (la table), "
                "« l' » devant une voyelle (l'ami, l'eau), « les » (les livres). Le genre du nom (masculin ou "
                "féminin) détermine l'article ; il faut l'apprendre avec chaque mot."},
    {"title": "Le verbe « être » au présent (A1)", "source_type": "grammar",
     "content": "Conjugaison du verbe être au présent : je suis, tu es, il/elle est, nous sommes, "
                "vous êtes, ils/elles sont. Exemples : « Je suis étudiant. », « Tu es fatigué ? », "
                "« Nous sommes à Paris. », « Ils sont contents. » C'est un verbe irrégulier très fréquent."},
    {"title": "Le verbe « avoir » au présent (A1)", "source_type": "grammar",
     "content": "Conjugaison du verbe avoir au présent : j'ai, tu as, il/elle a, nous avons, vous avez, "
                "ils/elles ont. Exemples : « J'ai vingt ans. », « Tu as un stylo ? », « Elle a faim. » "
                "En français, on utilise avoir pour l'âge (J'ai 20 ans) et pour la faim/soif (j'ai faim, j'ai soif)."},

    # --- A2 : quotidien, passé composé ---
    {"title": "Le passé composé (A2)", "source_type": "grammar",
     "content": "Le passé composé exprime une action terminée. Il se forme avec un auxiliaire (avoir ou être) "
                "au présent + le participe passé. Avec avoir : « J'ai mangé une pomme. » Avec être (verbes de "
                "mouvement et pronominaux) : « Je suis allé(e) au marché. » Avec être, le participe s'accorde avec "
                "le sujet : « Elle est partie. », « Ils sont arrivés. »"},
    {"title": "Au restaurant (A2)", "source_type": "dialogue",
     "content": "Serveur : « Bonjour, vous avez choisi ? » Client : « Oui, je voudrais le plat du jour, "
                "s'il vous plaît. » Serveur : « Et comme boisson ? » Client : « Une carafe d'eau et un verre de "
                "vin rouge. » Pour demander l'addition : « L'addition, s'il vous plaît. » Expression utile : "
                "« je voudrais » (conditionnel de politesse) est plus poli que « je veux »."},
    {"title": "Demander son chemin (A2)", "source_type": "dialogue",
     "content": "Pour demander : « Excusez-moi, où est la gare, s'il vous plaît ? » Réponses possibles : "
                "« Allez tout droit. », « Tournez à gauche / à droite. », « C'est à côté de la banque. », "
                "« C'est en face de l'église. », « Prenez la première rue à droite. » Vocabulaire : tout droit, "
                "à gauche, à droite, en face de, à côté de, près de, loin de."},

    # --- B1 : opinions, futur, imparfait ---
    {"title": "Imparfait vs passé composé (B1)", "source_type": "grammar",
     "content": "L'imparfait décrit un décor, une habitude ou une action en cours dans le passé : "
                "« Quand j'étais petit, je jouais au foot tous les jours. » Le passé composé exprime une action "
                "ponctuelle et terminée : « Hier, j'ai joué au foot. » Souvent on combine les deux : "
                "« Je dormais (imparfait, décor) quand le téléphone a sonné (passé composé, événement). »"},
    {"title": "Exprimer son opinion (B1)", "source_type": "lesson",
     "content": "Pour donner son avis : « Je pense que… », « Je crois que… », « À mon avis… », « Selon moi… ». "
                "Pour être d'accord : « Je suis d'accord. », « Tu as raison. », « Tout à fait. » Pour être en "
                "désaccord poliment : « Je ne suis pas d'accord. », « Je ne pense pas que ce soit vrai. » Notez : "
                "après « je ne pense pas que », on emploie souvent le subjonctif."},
    {"title": "Le futur simple (B1)", "source_type": "grammar",
     "content": "Le futur simple exprime une action future. Pour les verbes réguliers, on ajoute les terminaisons "
                "-ai, -as, -a, -ons, -ez, -ont à l'infinitif : « parler » → je parlerai, tu parleras, il parlera… "
                "Verbes irréguliers fréquents : être → je serai ; avoir → j'aurai ; aller → j'irai ; faire → je ferai. "
                "Exemple : « Demain, nous irons à la plage et nous ferons un pique-nique. »"},

    # --- B2 / C1 : nuances, subjonctif, registres ---
    {"title": "Le subjonctif présent (B2)", "source_type": "grammar",
     "content": "Le subjonctif exprime le doute, le souhait, l'émotion ou la nécessité, souvent après « que ». "
                "Il suit des expressions comme « il faut que », « je veux que », « bien que », « pour que ». "
                "Formation : radical de la 3e personne du pluriel au présent + -e, -es, -e, -ions, -iez, -ent. "
                "Exemple : « Il faut que tu fasses tes devoirs. », « Je suis content que tu sois là. » "
                "Verbes irréguliers : être (que je sois), avoir (que j'aie), aller (que j'aille), faire (que je fasse)."},
    {"title": "Connecteurs logiques (B2)", "source_type": "lesson",
     "content": "Pour structurer un texte ou un argument : la cause — « parce que », « car », « puisque », "
                "« étant donné que » ; la conséquence — « donc », « par conséquent », « c'est pourquoi » ; "
                "l'opposition — « mais », « cependant », « néanmoins », « en revanche » ; l'addition — « de plus », "
                "« en outre », « par ailleurs » ; la conclusion — « en somme », « finalement », « pour conclure »."},
    {"title": "Registres de langue (C1)", "source_type": "lesson",
     "content": "Le français distingue plusieurs registres. Familier : « Je suis crevé. », « C'est nul. », "
                "« bouquin » (livre), « bagnole » (voiture). Standard/courant : « Je suis fatigué. », « Ce n'est pas "
                "bien. » Soutenu/formel : « Je suis épuisé. », « Cela me semble regrettable. » Choisir le bon registre "
                "selon le contexte (ami, travail, écrit officiel) est une compétence clé au niveau avancé."},
    {"title": "Expressions idiomatiques courantes (C1)", "source_type": "lesson",
     "content": "Quelques expressions imagées : « avoir le cafard » = être triste ; « coûter les yeux de la tête » "
                "= être très cher ; « poser un lapin à quelqu'un » = ne pas venir à un rendez-vous ; « tomber dans "
                "les pommes » = s'évanouir ; « il pleut des cordes » = il pleut très fort ; « ce n'est pas la mer à "
                "boire » = ce n'est pas si difficile. Ces expressions ne se traduisent pas mot à mot."},
]

def main():
    os.makedirs(OUT_DIR, exist_ok=True)
    with open(OUT, "w", encoding="utf-8") as f:
        for d in DOCS:
            d.setdefault("audience", "french")
            d.setdefault("brand", "Les Petits Horizons")
            d.setdefault("language", "fr")
            d.setdefault("product", None)
            d["content_hash"] = hashlib.sha256(d["content"].encode()).hexdigest()
            f.write(json.dumps(d, ensure_ascii=False) + "\n")
    print(f"wrote {len(DOCS)} docs -> {OUT}")

if __name__ == "__main__":
    main()
