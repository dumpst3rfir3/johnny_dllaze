echo "Building goader"
GOOS=windows go build -o goader.exe
echo  `ls -al goader.exe`
echo `md5sum goader.exe`
echo "Done"
