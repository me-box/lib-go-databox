go get github.com/davecheney/godoc2md
echo "Re-generate README.md"
godoc2md ./ > README.md

ed README.md << HERE
$
d
d
w
HERE

echo '## Development of databox was supported by the following funding' >> README.md
echo '```' >> README.md
echo 'EP/N028260/1, Databox: Privacy-Aware Infrastructure for Managing Personal Data' >> README.md
echo 'EP/N028260/2, Databox: Privacy-Aware Infrastructure for Managing Personal Data' >> README.md
echo 'EP/N014243/1, Future Everyday Interaction with the Autonomous Internet of Things' >> README.md
echo 'EP/M001636/1, Privacy-by-Design: Building Accountability into the Internet of Things EP/M02315X/1, From Human Data to Personal Experience' >> README.md
echo '```' >> README.md
