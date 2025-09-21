# Structure 
Ce programme permet de se connecter à une base de données MongoDB et d’y insérer des données issues de fichiers XML.

# Fonctionnement
L’exécutable sert à exporter le contenu d’un fichier XML dans une base MongoDB.
Deux modes d’utilisation sont possibles :
- Lancer l’exécutable et indiquer le chemin d’accès du fichier XML.
- Glisser-déposer le fichier XML directement sur l’exécutable (drag & drop).

# Format attendu du fichier XML
Le fichier xml devra suivre cette structure : 
```xml
<AssociationsProfilsFonctionsEtSousFonctions>
    <AssociationProfilFonctionsEtSousFonctions>
      <Profil Code="D" Nom="" />
      <Fonctions>
        <Fonction Code="" Nom="" />
      </Fonctions>
      <SousFonctions>
        <SousFonction Code="" Nom="" />
      </SousFonctions>
    </AssociationProfilFonctionsEtSousFonctions>
    <AssociationProfilFonctionsEtSousFonctions>
      ...
    </AssociationProfilFonctionsEtSousFonctions>
    <AssociationProfilFonctionsEtSousFonctions>
      ...
    </AssociationProfilFonctionsEtSousFonctions>
</AssociationsProfilsFonctionsEtSousFonctions>
```
# Exemple 
## Fichier XML
```xml
<AssociationsProfilsFonctionsEtSousFonctions>
    <AssociationProfilFonctionsEtSousFonctions>
        <Profil Code="D01" Nom="Développeur" />
        <Fonctions>
            <Fonction Code="F01" Nom="Développement" />
            <Fonction Code="F02" Nom="Maintenance" />
        </Fonctions>
        <SousFonctions>
            <SousFonction Code="SF01" Nom="Front-End" />
            <SousFonction Code="SF02" Nom="Back-End" />
        </SousFonctions>
    </AssociationProfilFonctionsEtSousFonctions>
</AssociationsProfilsFonctionsEtSousFonctions>
```
## Exemple Insertion MongoDB

```json
{
  "Profil": { "Code": "D01", "Nom": "Développeur" },
  "Fonctions": [
    { "Code": "F01", "Nom": "Développement" },
    { "Code": "F02", "Nom": "Maintenance" }
  ],
  "SousFonctions": [
    { "Code": "SF01", "Nom": "Front-End" },
    { "Code": "SF02", "Nom": "Back-End" }
  ]
}
```
# Limitations
Le programme ne gère pas les doublons présents dans les sections :
- <Fonctions>
- <SousFonctions>
Ces doublons seront insérés tels quels dans la base.
