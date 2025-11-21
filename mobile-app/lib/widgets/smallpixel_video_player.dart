/// StreamVerse Video Player with Smallpixel Integration
/// Flutter widget for bandwidth-optimized video playback

import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';
import '../services/smallpixel_service.dart';

class SmallpixelVideoPlayer extends StatefulWidget {
  final String videoUrl;
  final bool autoPlay;
  final bool enableSmallpixel;
  final SmallpixelConfig? smallpixelConfig;
  final VoidCallback? onEnded;

  const SmallpixelVideoPlayer({
    Key? key,
    required this.videoUrl,
    this.autoPlay = false,
    this.enableSmallpixel = true,
    this.smallpixelConfig,
    this.onEnded,
  }) : super(key: key);

  @override
  State<SmallpixelVideoPlayer> createState() => _SmallpixelVideoPlayerState();
}

class _SmallpixelVideoPlayerState extends State<SmallpixelVideoPlayer> {
  late VideoPlayerController _controller;
  SmallpixelService? _smallpixelService;
  UpscalingStats? _currentStats;
  double _totalBandwidthSaved = 0.0;
  bool _showStats = false;

  @override
  void initState() {
    super.initState();
    _initializePlayer();
  }

  Future<void> _initializePlayer() async {
    // Initialize video player
    _controller = VideoPlayerController.network(widget.videoUrl);
    await _controller.initialize();

    // Initialize Smallpixel if enabled
    if (widget.enableSmallpixel) {
      final config = widget.smallpixelConfig ??
          SmallpixelConfig(
            apiKey: 'YOUR_SMALLPIXEL_API_KEY',
            targetResolution: TargetResolution.auto,
            quality: UpscalingQuality.high,
          );

      _smallpixelService = SmallpixelService(config);
      await _smallpixelService!.initialize(_controller);
      await _smallpixelService!.startUpscaling();

      // Listen to stats updates
      _smallpixelService!.statsStream?.listen((stats) {
        setState(() {
          _currentStats = stats;
          _totalBandwidthSaved = _smallpixelService!.totalBandwidthSaved;
        });
      });
    }

    // Auto-play if enabled
    if (widget.autoPlay) {
      _controller.play();
    }

    // Listen for video end
    _controller.addListener(() {
      if (_controller.value.position == _controller.value.duration) {
        widget.onEnded?.call();
      }
    });

    setState(() {});
  }

  @override
  void dispose() {
    _controller.dispose();
    _smallpixelService?.dispose();
    super.dispose();
  }

  void _toggleStats() {
    setState(() {
      _showStats = !_showStats;
    });
  }

  void _toggleSmallpixel() async {
    if (_smallpixelService?.isUpscaling ?? false) {
      await _smallpixelService!.stopUpscaling();
    } else {
      if (_smallpixelService != null) {
        await _smallpixelService!.startUpscaling();
      }
    }
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    if (!_controller.value.isInitialized) {
      return const Center(
        child: CircularProgressIndicator(),
      );
    }

    return Stack(
      alignment: Alignment.center,
      children: [
        // Video player
        AspectRatio(
          aspectRatio: _controller.value.aspectRatio,
          child: VideoPlayer(_controller),
        ),

        // Play/Pause button
        Positioned(
          child: GestureDetector(
            onTap: () {
              setState(() {
                _controller.value.isPlaying
                    ? _controller.pause()
                    : _controller.play();
              });
            },
            child: Icon(
              _controller.value.isPlaying ? Icons.pause : Icons.play_arrow,
              size: 80,
              color: Colors.white.withOpacity(0.7),
            ),
          ),
        ),

        // Smallpixel stats overlay
        if (_showStats && _currentStats != null)
          Positioned(
            top: 16,
            right: 16,
            child: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.black.withOpacity(0.8),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      const Icon(
                        Icons.arrow_upward,
                        color: Colors.greenAccent,
                        size: 16,
                      ),
                      const SizedBox(width: 4),
                      const Text(
                        'Smallpixel Active',
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          fontSize: 12,
                        ),
                      ),
                      const SizedBox(width: 8),
                      GestureDetector(
                        onTap: _toggleStats,
                        child: const Icon(
                          Icons.close,
                          color: Colors.white54,
                          size: 16,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  _buildStatRow('Source:', _currentStats!.originalResolution),
                  _buildStatRow('Target:', _currentStats!.targetResolution),
                  _buildStatRow(
                    'Latency:',
                    '${_currentStats!.upscalingLatencyMs.toStringAsFixed(1)}ms',
                  ),
                  _buildStatRow(
                    'Saved:',
                    '${_totalBandwidthSaved.toStringAsFixed(2)} MB',
                  ),
                  _buildStatRow(
                    'Savings:',
                    '${_currentStats!.savingsPercentage.toStringAsFixed(1)}%',
                  ),
                  _buildStatRow('FPS:', '${_currentStats!.frameRate.toInt()}'),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: _toggleSmallpixel,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.redAccent,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 6,
                        ),
                      ),
                      child: const Text(
                        'Disable',
                        style: TextStyle(fontSize: 11),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),

        // Smallpixel toggle button (when stats hidden)
        if (!_showStats && widget.enableSmallpixel)
          Positioned(
            top: 16,
            right: 16,
            child: GestureDetector(
              onTap: _toggleStats,
              child: Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 8,
                ),
                decoration: BoxDecoration(
                  color: (_smallpixelService?.isUpscaling ?? false)
                      ? Colors.greenAccent
                      : Colors.blueAccent,
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(
                      (_smallpixelService?.isUpscaling ?? false)
                          ? Icons.arrow_upward
                          : Icons.arrow_downward,
                      color: Colors.white,
                      size: 16,
                    ),
                    const SizedBox(width: 4),
                    Text(
                      (_smallpixelService?.isUpscaling ?? false)
                          ? 'Saving ${_totalBandwidthSaved.toStringAsFixed(0)} MB'
                          : 'Enable Smallpixel',
                      style: const TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 11,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),

        // Video controls
        Positioned(
          bottom: 0,
          left: 0,
          right: 0,
          child: Container(
            color: Colors.black.withOpacity(0.5),
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Row(
              children: [
                Text(
                  _formatDuration(_controller.value.position),
                  style: const TextStyle(color: Colors.white, fontSize: 12),
                ),
                Expanded(
                  child: Slider(
                    value: _controller.value.position.inSeconds.toDouble(),
                    max: _controller.value.duration.inSeconds.toDouble(),
                    onChanged: (value) {
                      _controller.seekTo(Duration(seconds: value.toInt()));
                    },
                    activeColor: Colors.greenAccent,
                  ),
                ),
                Text(
                  _formatDuration(_controller.value.duration),
                  style: const TextStyle(color: Colors.white, fontSize: 12),
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildStatRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 2),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: const TextStyle(
              color: Colors.white70,
              fontSize: 11,
            ),
          ),
          const SizedBox(width: 8),
          Text(
            value,
            style: const TextStyle(
              color: Colors.white,
              fontSize: 11,
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
    );
  }

  String _formatDuration(Duration duration) {
    String twoDigits(int n) => n.toString().padLeft(2, '0');
    String minutes = twoDigits(duration.inMinutes.remainder(60));
    String seconds = twoDigits(duration.inSeconds.remainder(60));
    return '$minutes:$seconds';
  }
}
