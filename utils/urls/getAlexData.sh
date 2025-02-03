#/usr/bin/bash
wget -q http://s3.amazonaws.com/alexa-static/top-1m.csv.zip;unzip top-1m.csv.zip; awk -F ',' '{print $2}' top-1m.csv.deprecated|head -100000 > alexa-top-100000.txt; rm top-1m.csv*