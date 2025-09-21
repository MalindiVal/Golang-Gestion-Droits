# Structure 
Ce programme permet de se connecter à une base de donnée mongoDB 
# Fonctionnement
Cette executable consiste à exporter des fichiers xml sur une base de donnée mongoDB :

Il suffit de lancer l'executable puis d'indiquer le chemin d'accées du fichier 
ou encore de prendre le fichier en maintenant le clique gauche et l'emmener directement au dessus de l'executable 

Le fichier xml devra suivre cette structure : 

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

Attention : Le programme ne prend pas en compte les doublons se trouvant dans la section "fonctions" et "sous fonctions" d'une association
