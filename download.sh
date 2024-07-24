#!/bin/bash

# URL de base pour les fichiers audio des sourates
base_url="https://download.quranicaudio.com/quran/abdulbaset_warsh/"

# Répertoire où les fichiers seront téléchargés
download_directory="quran_sourates"

# Crée le répertoire de téléchargement s'il n'existe pas
mkdir -p "$download_directory"

# Télécharger les 114 sourates
for i in $(seq -f "%03g" 1 114)
do
    file_url="${base_url}${i}.mp3"
    echo $file_url
    # wget -P "$download_directory" "$file_url"
done

echo "Téléchargement terminé."
