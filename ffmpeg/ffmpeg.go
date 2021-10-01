package ffmpeg

import (
	"github.com/pkg/errors"
)

// VideoInfo
/**
ffprobe version 4.2.4-1ubuntu0.1 Copyright (c) 2007-2020 the FFmpeg developers
  built with gcc 9 (Ubuntu 9.3.0-10ubuntu2)
  configuration: --prefix=/usr --extra-version=1ubuntu0.1 --toolchain=hardened --libdir=/usr/lib/x86_64-linux-gnu --incdir=/usr/include/x86_64-linux-gnu --arch=amd64 --enable-gpl --disable-stripping --enable-avresample --disable-filter=resample --enable-avisynth --enable-gnutls --enable-ladspa --enable-libaom --enable-libass --enable-libbluray --enable-libbs2b --enable-libcaca --enable-libcdio --enable-libcodec2 --enable-libflite --enable-libfontconfig --enable-libfreetype --enable-libfribidi --enable-libgme --enable-libgsm --enable-libjack --enable-libmp3lame --enable-libmysofa --enable-libopenjpeg --enable-libopenmpt --enable-libopus --enable-libpulse --enable-librsvg --enable-librubberband --enable-libshine --enable-libsnappy --enable-libsoxr --enable-libspeex --enable-libssh --enable-libtheora --enable-libtwolame --enable-libvidstab --enable-libvorbis --enable-libvpx --enable-libwavpack --enable-libwebp --enable-libx265 --enable-libxml2 --enable-libxvid --enable-libzmq --enable-libzvbi --enable-lv2 --enable-omx --enable-openal --enable-opencl --enable-opengl --enable-sdl2 --enable-libdc1394 --enable-libdrm --enable-libiec61883 --enable-nvenc --enable-chromaprint --enable-frei0r --enable-libx264 --enable-shared
  libavutil      56. 31.100 / 56. 31.100
  libavcodec     58. 54.100 / 58. 54.100
  libavformat    58. 29.100 / 58. 29.100
  libavdevice    58.  8.100 / 58.  8.100
  libavfilter     7. 57.100 /  7. 57.100
  libavresample   4.  0.  0 /  4.  0.  0
  libswscale      5.  5.100 /  5.  5.100
  libswresample   3.  5.100 /  3.  5.100
  libpostproc    55.  5.100 / 55.  5.100
Input #0, matroska,webm, from 'The.Prince.of.Egypt.1998.mkv':
  Metadata:
    title           : Despicable Me - YIFY
    encoder         : libebml v1.3.0 + libmatroska v1.4.0
    creation_time   : 2015-06-27T07:42:32.000000Z
  Duration: 01:34:40.34, start: 0.000000, bitrate: 1578 kb/s
    Stream #0:0(eng): Video: h264 (High), yuv420p(progressive), 1280x720 [SAR 1:1 DAR 16:9], 23.98 fps, 23.98 tbr, 1k tbn, 47.95 tbc (default)
    Metadata:
      title           : The.Prince.of.Egypt.1998-EN580
    Stream #0:1(eng): Audio: aac (LC), 48000 Hz, 5.1, fltp (default)
    Metadata:
      title           : The.Prince.of.Egypt.1998-EN580
    Stream #0:2(eng): Subtitle: subrip (default)
    Metadata:
      title           : 英文字幕『en580.com』
    Stream #0:3(chi): Subtitle: subrip
    Metadata:
      title           : 中上英下『en580.com』
    Stream #0:4(eng): Subtitle: subrip
    Metadata:
      title           : 英上中下『en580.com』
    Stream #0:5(chi): Subtitle: subrip
    Metadata:
      title           : 中文字幕『en580.com』
*/
type VideoInfo struct {
	Streams StreamList
}

type StreamType string

const (
	StreamTypeVideo    = StreamType("Video")
	StreamTypeAudio    = StreamType("Audio")
	StreamTypeSubtitle = StreamType("Subtitle")
)

type StreamList []Stream

type Stream struct {
	Index              int          `json:"index"`
	CodecName          string       `json:"codec_name"`
	CodecLongName      string       `json:"codec_long_name"`
	Profile            string       `json:"profile"`
	CodecType          string       `json:"codec_type"`
	CodecTimeBase      string       `json:"codec_time_base"`
	CodecTagString     string       `json:"codec_tag_string"`
	CodecTag           string       `json:"codec_tag"`
	Width              int64        `json:"width"`
	Height             int64        `json:"height"`
	CodedWidth         int64        `json:"coded_width"`
	CodedHeight        int64        `json:"coded_height"`
	HasBFrames         IntBool      `json:"has_b_frames"`
	SampleAspectRatio  string       `json:"sample_aspect_ratio"`
	DisplayAspectRatio string       `json:"display_aspect_ratio"`
	PixFmt             string       `json:"pix_fmt"`
	Level              int          `json:"level"`
	ChromaLocation     string       `json:"chroma_location"`
	FieldOrder         string       `json:"field_order"`
	Refs               int          `json:"refs"`
	IsAvc              StringBool   `json:"is_avc"`
	NalLengthSize      StringNumber `json:"nal_length_size"`
	RFrameRate         string       `json:"r_frame_rate"`
	AvgFrameRate       string       `json:"avg_frame_rate"`
	TimeBase           string       `json:"time_base"`
	StartPts           int          `json:"start_pts"`
	StartTime          StringFloat  `json:"start_time"`
	BitsPerRawSample   StringNumber `json:"bits_per_raw_sample"`
	Disposition        Disposition  `json:"disposition"`
	Tags               Tags         `json:"tags"`
}

type Tags map[string]string

type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

func GetVideoInfo(fileName string) (videoInfo *VideoInfo, err error) {
	data, err := FFProbe(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "get video info by ffprobe failed")
	}
	return parseData(data)
}
