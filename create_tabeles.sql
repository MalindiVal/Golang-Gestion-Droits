CREATE TABLE Profil(
   CodeProfil VARCHAR(50),
   Nom VARCHAR(50),
   PRIMARY KEY(CodeProfil)
);

CREATE TABLE Fonctions(
   CodeFonction VARCHAR(50),
   Nom VARCHAR(50) NOT NULL,
   PRIMARY KEY(CodeFonction)
);

CREATE TABLE Assosiations(
   Id VARCHAR(50),
   Type VARCHAR(50),
   CodeFonction VARCHAR(50) NOT NULL,
   CodeProfil VARCHAR(50) NOT NULL,
   PRIMARY KEY(Id),
   FOREIGN KEY(CodeFonction) REFERENCES Fonctions(CodeFonction),
   FOREIGN KEY(CodeProfil) REFERENCES Profil(CodeProfil)
);
