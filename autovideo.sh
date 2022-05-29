#!/bin/bash

INPUT_AUDIO_DIR=$(realpath "$1")
INPUT_VIDEO_DIR=$(realpath "$2")
OUTPUT_FILE=$(realpath "$3")
cd "$(dirname "$0")" || exit

function printStart {
  echo "#######################################"
  echo ""
  echo ""
}

function printEnd {
  echo ""
  echo ""
  echo "#######################################"
  sleep 2
}

function printLine {
  echo "#######################################"
  echo ""
  echo ""
  echo "$1"
  echo ""
  echo ""
  echo "#######################################"
  sleep 1
}

# $1: dir
function lsOrder {
  ls -tr $1 -1 || exit
}

# $1: dir
# $2: outDir
function resizeAll {
  mkdir -p "$2" || exit

  lsOrder "$1" | while read f; do
    ffmpeg -nostdin -i "$1/$f" -y \
      -vf "scale=w=1280:h=720:force_original_aspect_ratio=1,pad=1280:720:(ow-iw)/2:(oh-ih) /2" \
      -vf "fps=24" \
      -vcodec libx264 -an "$2/$f"
  done
}

# $1: dir
# $2: target
function concatVideos {
  COUNT=$(lsOrder "$1" | wc -l)

  if [ "$COUNT" -eq "0" ]; then
    printLine "No videos"
    return
  fi

  TEMP=temp.concatVideos.txt
  rm $TEMP
  lsOrder "$1" | while read f; do
    echo "file 'c:${1:2}/$f'" >>$TEMP
  done

  ffmpeg -f concat -safe 0 -i $TEMP -vcodec libx264 -an -y "$2" || exit
  rm $TEMP
}

# $1: dir
# $2: output
function concatWaves {
  TEMP=temp.concatWaves.txt
  lsOrder "$1" | grep ".*\.wav" |
    while read f; do
      echo "file 'c:${1:2}/$f'" >>$TEMP
    done

  ffmpeg -f concat -safe 0 -i $TEMP -y "$2"
  rm $TEMP
}

# $1: file
function getLengthSeconds {
  ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$1"
}

# $1: str
# $2: str
function ensureGt {
  AWK="BEGIN {print ("$1" > "$2")}"
  BOOL="$(awk "$AWK")"
  if [ "$BOOL" -eq "0" ]; then
    echo "$1 <= $2"
    exit 1
  fi
}

# $1: video
# $2: audio
# $3: output
function merge {
  ffmpeg -i "$1" -i "$2" -c:v copy -c:a aac -shortest "$3" || exit
}

TEMP_VIDEO=temp.video.mp4
TEMP_AUDIO=temp.audio.mp3
RESIZED_VIDEO_DIR=$(realpath videos_resized)

rm "$TEMP_VIDEO"
rm "$TEMP_AUDIO"
rm -rf "$RESIZED_VIDEO_DIR"

printStart
echo "input audio directory: $INPUT_AUDIO_DIR"
echo "input video directory: $INPUT_VIDEO_DIR"
echo "output file: $OUTPUT_FILE"
printEnd

resizeAll "$INPUT_VIDEO_DIR" "$RESIZED_VIDEO_DIR" || exit
concatVideos "$RESIZED_VIDEO_DIR" "$TEMP_VIDEO" || exit
concatWaves "$INPUT_AUDIO_DIR" "$TEMP_AUDIO" || exit
AUDIO_LENGTH=$(getLengthSeconds $TEMP_AUDIO)
VIDEO_LENGTH=$(getLengthSeconds "$TEMP_VIDEO")

printStart
echo "audio length: $AUDIO_LENGTH"
echo "video length: $VIDEO_LENGTH"
printEnd

ensureGt "$VIDEO_LENGTH" "$AUDIO_LENGTH"
merge "$TEMP_VIDEO" "$TEMP_AUDIO" "$OUTPUT_FILE"

printLine "finish"
