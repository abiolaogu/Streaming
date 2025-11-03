/**
 * StreamVerse Flutter Video Player SDK
 * Issue #28: Video Player SDK Development
 */

import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';
import 'dart:async';

class DRMConfig {
  final String? widevineLicenseUrl;
  final String? fairplayLicenseUrl;
  final String? playreadyLicenseUrl;
  final String? certificateUrl;

  DRMConfig({
    this.widevineLicenseUrl,
    this.fairplayLicenseUrl,
    this.playreadyLicenseUrl,
    this.certificateUrl,
  });
}

class PlayerConfig {
  final String videoUrl;
  final DRMConfig? drm;
  final bool abrEnabled;
  final int bufferingGoal; // seconds
  final Function(String)? onQualityChange;
  final Function(String)? onError;
  final Function()? onPlay;
  final Function()? onPause;
  final Function()? onEnded;

  PlayerConfig({
    required this.videoUrl,
    this.drm,
    this.abrEnabled = true,
    this.bufferingGoal = 8,
    this.onQualityChange,
    this.onError,
    this.onPlay,
    this.onPause,
    this.onEnded,
  });
}

class StreamVersePlayer extends StatefulWidget {
  final PlayerConfig config;

  const StreamVersePlayer({
    Key? key,
    required this.config,
  }) : super(key: key);

  @override
  State<StreamVersePlayer> createState() => _StreamVersePlayerState();
}

class _StreamVersePlayerState extends State<StreamVersePlayer> {
  VideoPlayerController? _controller;
  bool _isInitialized = false;
  bool _isPlaying = false;
  String? _currentQuality;
  List<String> _availableQualities = [];

  @override
  void initState() {
    super.initState();
    _initializePlayer();
  }

  Future<void> _initializePlayer() async {
    try {
      _controller = VideoPlayerController.networkUrl(
        Uri.parse(widget.config.videoUrl),
        videoPlayerOptions: VideoPlayerOptions(
          mixWithOthers: true,
          allowBackgroundPlayback: false,
        ),
      );

      await _controller!.initialize();
      
      _controller!.addListener(() {
        if (_controller!.value.isPlaying && !_isPlaying) {
          _isPlaying = true;
          widget.config.onPlay?.call();
        } else if (!_controller!.value.isPlaying && _isPlaying) {
          _isPlaying = false;
          widget.config.onPause?.call();
        }

        if (_controller!.value.position >= _controller!.value.duration &&
            _controller!.value.duration.inMilliseconds > 0) {
          widget.config.onEnded?.call();
        }
      });

      setState(() {
        _isInitialized = true;
      });
    } catch (e) {
      widget.config.onError?.call(e.toString());
    }
  }

  void play() {
    _controller?.play();
  }

  void pause() {
    _controller?.pause();
  }

  void seekTo(Duration position) {
    _controller?.seekTo(position);
  }

  void setQuality(String qualityId) {
    // TODO: Implement quality switching
    _currentQuality = qualityId;
    widget.config.onQualityChange?.call(qualityId);
  }

  @override
  void dispose() {
    _controller?.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (!_isInitialized) {
      return const Center(child: CircularProgressIndicator());
    }

    return AspectRatio(
      aspectRatio: _controller!.value.aspectRatio,
      child: VideoPlayer(_controller!),
    );
  }
}

