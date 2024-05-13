
/*
* XOR cipher
* XOR Encryption is an encryption method used to encrypt data and is 
* hard to crack by brute-force method, i.e generating random encryption 
* keys to match with the correct one. 
* 
* Code by @Anvikshajain
*/

//Code to encrypt data device and calculate the hast to send to the gateway application

#include<string.h>

//#include "xor_chiper.h"
#include <stdio.h>
#include <stdlib.h>


//Get digest of data device
size_t getHash(const char* cp)
{
    size_t hash = 0;
    while (*cp)
        hash = (hash * 10) + *cp++ - '0';
    return hash;
}

// The same function is used to encrypt and
// decrypt data device
void encryptDecrypt(char inpString[],char key)
{ 
	// Define XOR key 
	// Any character value will work 
	char xorKey = key;
	// calculate length of input string 
	int len = strlen(inpString);
	// perform XOR operation of key 
	// with every caracter in string 
	for (int i = 0; i < len; i++) 
	{
	    if(i == len-1){
	        continue;
	    }

		inpString[i] = inpString[i] ^ xorKey; 

	}
}

int main(int argc, char **argv){

    FILE *fpR; //File to receive clean data
    FILE *fpW; //File to put encrypted data

    char *line = NULL;
    size_t len = 0;
    ssize_t read;

    if (argc < 4)
        return 1;

    fpR = fopen(argv[1], "r");
    fpW = fopen(argv[2], "a");
    if (fpR == NULL)
        exit(EXIT_FAILURE);
    if (fpW == NULL)
        exit(EXIT_FAILURE);

   while ((read = getline(&line, &len, fpR)) != -1) {
        if(strcmp(argv[3],"in") == 0){ //Encrypt data
             size_t hashl = getHash(line);
             printf("Retrieved line of length %zu and hash %zu :\n", read,hashl);
             encryptDecrypt(line,argv[4]);
             fprintf(fpW, "%s", line);
             fprintf(fpW, "%zu\n", hashl);
        }
        if (strcmp(argv[3],"de") == 0){ //Decrypt data
            printf("Retrieved line of length %zu :\n", read);
            encryptDecrypt(line,argv[4]);
            fprintf(fpW, "%s", line);
        }
   }

   free(line);

   fclose(fpW);
   fclose(fpR);


   return 0;


}



// Driver program to test above function 
/*int main()
{
	char originalString[] = "Technology";

	// Encrypt the string
	printf("Encrypted String: ");
	encryptDecrypt(originalString);
	printf("\n");

	// Decrypt the string
	printf("Decrypted String: ");
	encryptDecrypt(originalString);

	return 0;
} */
