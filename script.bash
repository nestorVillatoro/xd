mkdisk -size=5 -unit=M -fit=WF -path="/home/nestor-villatoro/go/src/MIA_2S_P1_202200252/discos/DiscoLab.mia"
fdisk -size=1 -type=P -unit=M -fit=BF -name="Particion1" -path="/home/nestor-villatoro/go/src/MIA_2S_P1_202200252/discos/DiscoLab.mia"
mount -name="Particion1" -path="/home/nestor-villatoro/go/src/MIA_2S_P1_202200252/discos/DiscoLab.mia"
mkfs -id=520A 
mkdir -path="/home"
mkdir -path="/home/usac"
mkdir -path="/home/work"
mkdir -path="/home/usac/mia"
mkfile -size=68 -path=/home/usac/mia/a.txt
rep -id=520A -path="/home/nestor-villatoro/go/src/MIA_2S_P1_202200252/output/report_mbr.png" -name=mbr