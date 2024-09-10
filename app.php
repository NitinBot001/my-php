<?php
// Define a function to fetch YouTube video details using yt-dlp
function getYouTubeVideoDetails($videoUrl) {
    // Define the path to yt-dlp executable
    $ytDlpPath = '/path/to/yt-dlp'; // Change this to the path of your yt-dlp installation

    // Run yt-dlp to get video info in JSON format
    $command = escapeshellcmd("$ytDlpPath -J --no-warnings --no-progress $videoUrl");
    $output = shell_exec($command);

    // Decode JSON output
    $videoData = json_decode($output, true);

    // Check if decoding was successful
    if ($videoData === null) {
        return array('error' => 'Failed to fetch video details.');
    }

    // Extract required details
    $title = isset($videoData['title']) ? $videoData['title'] : 'No title available';
    $thumbnailUrl = isset($videoData['thumbnail']) ? $videoData['thumbnail'] : 'No thumbnail available';
    $formats = isset($videoData['formats']) ? $videoData['formats'] : array();

    $audioUrl = 'No audio link available';
    foreach ($formats as $format) {
        if ($format['acodec'] !== 'none' && isset($format['url'])) {
            $audioUrl = $format['url'];
            break;
        }
    }

    // Return details in an associative array
    return array(
        'title' => $title,
        'thumbnailUrl' => $thumbnailUrl,
        'audioUrl' => $audioUrl
    );
}

// Example usage
if (isset($_GET['video_url'])) {
    $videoUrl = $_GET['video_url'];
    $videoDetails = getYouTubeVideoDetails($videoUrl);
    header('Content-Type: application/json');
    echo json_encode($videoDetails);
} else {
    echo json_encode(array('error' => 'No video URL provided.'));
}
?>
