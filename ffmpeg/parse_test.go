package ffmpeg

import "testing"

func TestParseData(t *testing.T) {
	input := []byte(`{
    "streams": [
        {
            "index": 0,
            "codec_name": "h264",
            "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
            "profile": "High",
            "codec_type": "video",
            "codec_time_base": "1001/48000",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "width": 1280,
            "height": 720,
            "coded_width": 1280,
            "coded_height": 720,
            "has_b_frames": 1,
            "sample_aspect_ratio": "1:1",
            "display_aspect_ratio": "16:9",
            "pix_fmt": "yuv420p",
            "level": 41,
            "chroma_location": "left",
            "field_order": "progressive",
            "refs": 1,
            "is_avc": "true",
            "nal_length_size": "4",
            "r_frame_rate": "24000/1001",
            "avg_frame_rate": "24000/1001",
            "time_base": "1/1000",
            "start_pts": 0,
            "start_time": "0.000000",
            "bits_per_raw_sample": "8",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "eng",
                "title": "The.Prince.of.Egypt.1998-EN580"
            }
        },
        {
            "index": 1,
            "codec_name": "aac",
            "codec_long_name": "AAC (Advanced Audio Coding)",
            "profile": "LC",
            "codec_type": "audio",
            "codec_time_base": "1/48000",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "sample_fmt": "fltp",
            "sample_rate": "48000",
            "channels": 6,
            "channel_layout": "5.1",
            "bits_per_sample": 0,
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/1000",
            "start_pts": 0,
            "start_time": "0.000000",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "eng",
                "title": "The.Prince.of.Egypt.1998-EN580"
            }
        },
        {
            "index": 2,
            "codec_name": "subrip",
            "codec_long_name": "SubRip subtitle",
            "codec_type": "subtitle",
            "codec_time_base": "0/1",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/1000",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 5680342,
            "duration": "5680.342000",
            "disposition": {
                "default": 1,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "eng",
                "title": "英文字幕『en580.com』"
            }
        },
        {
            "index": 3,
            "codec_name": "subrip",
            "codec_long_name": "SubRip subtitle",
            "codec_type": "subtitle",
            "codec_time_base": "0/1",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/1000",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 5680342,
            "duration": "5680.342000",
            "disposition": {
                "default": 0,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "chi",
                "title": "中上英下『en580.com』"
            }
        },
        {
            "index": 4,
            "codec_name": "subrip",
            "codec_long_name": "SubRip subtitle",
            "codec_type": "subtitle",
            "codec_time_base": "0/1",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/1000",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 5680342,
            "duration": "5680.342000",
            "disposition": {
                "default": 0,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "eng",
                "title": "英上中下『en580.com』"
            }
        },
        {
            "index": 5,
            "codec_name": "subrip",
            "codec_long_name": "SubRip subtitle",
            "codec_type": "subtitle",
            "codec_time_base": "0/1",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/1000",
            "start_pts": 0,
            "start_time": "0.000000",
            "duration_ts": 5680342,
            "duration": "5680.342000",
            "disposition": {
                "default": 0,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0,
                "timed_thumbnails": 0
            },
            "tags": {
                "language": "chi",
                "title": "中文字幕『en580.com』"
            }
        }
    ],
    "format": {
        "filename": "The.Prince.of.Egypt.1998.mkv",
        "nb_streams": 6,
        "nb_programs": 0,
        "format_name": "matroska,webm",
        "format_long_name": "Matroska / WebM",
        "start_time": "0.000000",
        "duration": "5680.342000",
        "size": "1121092516",
        "bit_rate": "1578908",
        "probe_score": 100,
        "tags": {
            "title": "Despicable Me - YIFY",
            "encoder": "libebml v1.3.0 + libmatroska v1.4.0",
            "creation_time": "2015-06-27T07:42:32.000000Z"
        }
    }
}`)
	video, err := parseData(input)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(video)
}
