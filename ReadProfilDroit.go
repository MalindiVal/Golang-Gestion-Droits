package main

//Importation des fonctions
import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"

	// Fonctions MongoDB
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Collection
var (
	usersCollection *mongo.Collection
	funcsCollection *mongo.Collection
	AssoCollection  *mongo.Collection
)

type Query struct {
	Associations []AssociationProfilFonctionsEtSousFonctions `xml:"AssociationsProfilsFonctionsEtSousFonctions>AssociationProfilFonctionsEtSousFonctions"`
}

// Type correspondant à une association
type AssociationProfilFonctionsEtSousFonctions struct {
	XMLName       xml.Name  `xml:"AssociationProfilFonctionsEtSousFonctions"`
	Profil        Element   `xml:"Profil"`
	Fonctions     []Element `xml:"Fonctions>Fonction"`
	SousFonctions []Element `xml:"SousFonctions>SousFonction"`
}

// Type correspondant à un element
type Element struct {
	Code string `xml:"Code,attr"`
	Nom  string `xml:"Nom,attr"`
}

// Connexion à la basse de données
func initClient() {

	fmt.Println("Connexion à la base de données ...")
	// Connection à la base de données
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017")) //
	if err != nil {
		fmt.Println("Echec de la connexion !!!")
		panic(err)
	} else {
		fmt.Println("Connexion réussie !!!")
	}

	// En cas de problème
	fmt.Println("Ping de la connexion ...")
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		fmt.Println("Test raté !!!")
		panic(err)
	} else {
		fmt.Println("Test réussi !!!")
	}

	fmt.Println("Initialisation de la liste des profils ...")
	usersCollection = client.Database("malindi").Collection("Profils")
	fmt.Println("Initialisation de la liste des fonctions ...")
	funcsCollection = client.Database("malindi").Collection("Fonctions")
	fmt.Println("Initialisation de la liste des associations ...")
	AssoCollection = client.Database("malindi").Collection("Associations")
}

// Permet de lire l fichier XML inséré
func readXml() (Query, error) {
	var q Query
	var fichier string
	if len(os.Args) > 1 {
		fichier = os.Args[1]
	} else {
		fmt.Println("Veuillez indiquer le nom de votre fichier xml :")
		n, err := fmt.Scanln(&fichier)
		if err != nil {
			fmt.Println("Erreur lors de la lecture:", err)
			return q, err
		}
		fmt.Println(n)
	}

	xmlFile, err := os.Open(fichier)
	//Si une errreur a été envoyé
	if err != nil {
		fmt.Println("Echec de l'ouverture du fichier !!!")
		fmt.Println("Error opening file:", err)
		return q, err
	} else {
		fmt.Println("Ouverture réussie !!!")
	}
	defer xmlFile.Close()

	// Verification du fichier en argument
	fmt.Println("Ouverture du fichier ...")

	//Lecture du fichier xml en argument
	fmt.Println("Lecture du fichier ...")
	b, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return q, err
	} else {
		fmt.Println("Lecture términée !!!")
	}

	fmt.Println("Convertion du fichier  ...")
	err = xml.Unmarshal(b, &q)
	if err != nil {
		fmt.Println("Error unmarshalling XML:", err)
		return q, err
	} else {
		fmt.Println("Convertion términée !!!")
	}

	return q, err
}

// Fontion Principal
func main() {

	initClient()

	q, err := readXml()
	if err != nil {
		fmt.Println("Echec de la conversion !!!")
	} else {
		fmt.Println("Réusssite de la conversion !!!")
		fmt.Println("Verification du format ...")
		if q.Associations != nil {
			fmt.Println("Format valide !!!")
			insertXml(deleteduplicate(q))
			fmt.Println("Transfert terminée !!!!!")
		} else {
			fmt.Println("Format XML invalide !!!")
			fmt.Println("Fin du programme !!!")
			return
		}

	}

	n, err := fmt.Scanln()
	if err == nil {
		fmt.Println(n)
		return
	}

}

// Insertion du fichier XML dans la base de données
func insertXml(q Query) {
	//Ajout des informations issues du fichier

	for _, association := range q.Associations {
		user := bson.D{{Key: "nom", Value: association.Profil.Nom}, {Key: "code", Value: association.Profil.Code}}

		var p Element

		//Vérification de la presence de boublons
		fmt.Printf("Verification de la presence du profil %s ...\n", association.Profil.Nom)
		//On verifie si , dans la base de données , le profil y est inscrit
		err := usersCollection.FindOne(context.TODO(), user).Decode(&p)
		//Si une errreur est détecté
		if err != nil {
			//Si cette erreur concerne le fait que la fonction dit que lke document rechercher n'existe pas
			if err == mongo.ErrNoDocuments {
				fmt.Println("Aucun doublon trouvé ")
				fmt.Printf("Insertion du profil %s ...\n", association.Profil.Nom)
				//Alors on réalisera une insertion
				usersCollection.InsertOne(context.TODO(), user)

				fmt.Printf("Ajout du profil %s réussi !!\n", association.Profil.Nom)
			}
		} else {
			fmt.Println("Doublon trouvé !!!")
		}

		fonctions := bson.A{}
		sousfonctions := bson.A{}

		//Verification de la présence de fonctions
		if len(association.Fonctions) != 0 {

			//Boucle d'ajout des fonctions dans la base de données
			for _, fonction := range association.Fonctions {

				user := bson.D{{Key: "nom", Value: fonction.Nom}, {Key: "code", Value: fonction.Code}}

				var fl Element
				// Supression des doublons
				fmt.Printf("Verification de la presence de la fonction %s ...\n", fonction.Nom)
				err := funcsCollection.FindOne(context.TODO(), user).Decode(&fl)
				if err != nil {
					// Aucun doublon trouvé
					if err == mongo.ErrNoDocuments {
						fmt.Println("Aucun doublon trouvé ")
						fmt.Printf("Insertion de la fonction %s ...\n", fonction.Nom)
						funcsCollection.InsertOne(context.TODO(), user)
						fmt.Printf("Ajout de la fonction %s réussi !!\n", fonction.Nom)
					}
				} else {
					fmt.Println("Doublon trouvé !!!")
				}

				fonctions = append(fonctions, user)

			}
		}

		//Verification de la présence de sous-fonction
		if len(association.SousFonctions) != 0 {

			for _, fonction := range association.SousFonctions {
				user := bson.D{{Key: "nom", Value: fonction.Nom}, {Key: "code", Value: fonction.Code}}

				var fl Element
				// Supression des doublons
				fmt.Printf("Verification de la presence de la sous-fonction %s ...\n", fonction.Nom)
				err := funcsCollection.FindOne(context.TODO(), user).Decode(&fl)
				if err != nil {
					// Aucun doublon trouvé
					if err == mongo.ErrNoDocuments {
						fmt.Println("Aucun doublon trouvé ")
						fmt.Printf("Insertion de la sous-fonction %s ...\n", fonction.Nom)
						funcsCollection.InsertOne(context.TODO(), user)
						fmt.Printf("Ajout de la fonction %s réussi !!\n", fonction.Nom)
					}
				} else {
					fmt.Println("Doublon trouvé !!!")
				}

				sousfonctions = append(sousfonctions, user)
			}

		}

		// Création de l'association
		//Suppression des doublons
		var a1 AssociationProfilFonctionsEtSousFonctions
		var a2 AssociationProfilFonctionsEtSousFonctions

		as := bson.D{{Key: "nomprofil", Value: association.Profil.Nom}}

		fmt.Printf("Verification de la presence de l'association ...\n")
		err = AssoCollection.FindOne(context.TODO(), as).Decode(&a1)
		if err == nil {
			asso := bson.D{{Key: "nomprofil", Value: association.Profil.Nom}, {Key: "codeprofil", Value: association.Profil.Code}, {Key: "fonctions", Value: fonctions}, {Key: "sousfonctions", Value: sousfonctions}}
			err = AssoCollection.FindOne(context.TODO(), asso).Decode(&a2)
			if err == nil {
				//fmt.Println("Doublon trouvé !!!")
			} else {
				if err == mongo.ErrNoDocuments {
					UpdateDoc(a1, association)
				}
			}
		} else {
			if err == mongo.ErrNoDocuments {
				fmt.Printf("Verification de la presence de l'association ...\n")
				asso := bson.D{{Key: "nomprofil", Value: association.Profil.Nom}, {Key: "codeprofil", Value: association.Profil.Code}, {Key: "fonctions", Value: fonctions}, {Key: "sousfonctions", Value: sousfonctions}}
				err = AssoCollection.FindOne(context.TODO(), asso).Decode(&a1)
				if err != nil {
					if err == mongo.ErrNoDocuments {
						fmt.Println("Aucun doublon trouvé ")
						fmt.Printf("Insertion de l'association ...\n")
						AssoCollection.InsertOne(context.TODO(), asso)

						fmt.Printf("Ajout de l'association réussi !!\n")
					}
				} else {
					fmt.Println("Doublon trouvé !!!")
					//UpdateDoc(a)
				}
			}
		}

	}

}

func remove(s []Element, i int) []Element {
	if i < 0 || i >= len(s) {
		// index out of range, return the original slice
		return s
	}
	s[i] = s[len(s)-1]
	// Ensure that the slice is not empty before truncating it
	if len(s) > 0 {
		return s[:len(s)-1]
	}
	return s
}

func deleteduplicate(q Query) Query {
	for _, association := range q.Associations {

		if len(association.Fonctions) != 0 {
			// Remove duplicates from Fonctions
			for j, fonction := range association.Fonctions {
				for i, search := range association.Fonctions {
					if fonction.Nom == search.Nom && i != j {
						association.Fonctions = remove(association.Fonctions, i)
					}
				}
			}
		}

		if len(association.SousFonctions) != 0 {
			// Remove duplicates from SousFonctions
			for j, sousFonction := range association.SousFonctions {
				for i, search := range association.SousFonctions {
					if sousFonction.Nom == search.Nom && i != j {
						association.SousFonctions = remove(association.SousFonctions, i)
					}
				}
			}
		}
	}
	return q
}

func UpdateDoc(a1 AssociationProfilFonctionsEtSousFonctions, a2 AssociationProfilFonctionsEtSousFonctions) {
	// Create Base.xml
	if err := writeXMLToFile("Base.xml", a1); err != nil {
		fmt.Println("Error creating/updating Base.xml: ", err)
		return
	}

	// Create New.xml
	if err := writeXMLToFile("New.xml", a2); err != nil {
		fmt.Println("Error creating/updating New.xml: ", err)
		return
	}

	// Prompt user for confirmation
	fmt.Println("Des modifications ont été réalisées sur l'association du profil ", a2.Profil.Nom)
	fmt.Println("Voulez-vous mettre à jour les données existantes avec les nouvelles ? Y/N")
	var res string
	fmt.Scanln(&res)
	if res == "Y" || res == "y" {
		updateDatabase(a2)
	}
}

// Génère un fichier XML
func writeXMLToFile(filename string, data interface{}) error {
	xmlFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	xmlFile.WriteString(xml.Header)
	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "\t")
	if err := encoder.Encode(data); err != nil {
		return err
	}

	return nil
}

func updateDatabase(a2 AssociationProfilFonctionsEtSousFonctions) {
	filter := bson.D{{Key: "nomprofil", Value: a2.Profil.Nom}}
	replacement := bson.D{
		{Key: "nomprofil", Value: a2.Profil.Nom},
		{Key: "codeprofil", Value: a2.Profil.Code},
		{Key: "fonctions", Value: a2.Fonctions},
		{Key: "sousfonctions", Value: a2.SousFonctions}, // Use underscore instead of space in field name
	}

	result, err := AssoCollection.ReplaceOne(context.TODO(), filter, replacement)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Mise à jour effectuée ", result)
	}
}
