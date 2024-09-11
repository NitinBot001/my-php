from flask import Flask, request, jsonify
import yt_dlp

app = Flask(__name__)

# Define a function to extract video metadata using yt-dlp
def get_video_metadata(search_query, username=None, password=None):
    ydl_opts = {
        'quiet': True,
        'noplaylist': True,
        'extract_flat': 'in_playlist',  # Only extract metadata, don't download
        'format': 'best'
    }

    # Add login credentials if provided
    if username and password:
        ydl_opts.update({
            'username': username,
            'password': password
        })

    with yt_dlp.YoutubeDL(ydl_opts) as ydl:
        try:
            # Use yt-dlp to search and get metadata
            result = ydl.extract_info(f"ytsearch:{search_query}", download=False)
            
            if 'entries' in result:
                # Return metadata of the first search result
                video = result['entries'][0]
                return {
                    'title': video.get('title'),
                    'id': video.get('id'),
                    'url': video.get('webpage_url'),
                    'duration': video.get('duration'),
                    'uploader': video.get('uploader'),
                    'view_count': video.get('view_count'),
                    'like_count': video.get('like_count'),
                    'description': video.get('description')
                }
            else:
                return {'error': 'No videos found'}

        except Exception as e:
            return {'error': str(e)}

# Define a route for video search and metadata fetching
@app.route('/search', methods=['GET'])
def search_video():
    query = request.args.get('query')
    username = request.args.get('username')
    password = request.args.get('password')
    
    if not query:
        return jsonify({'error': 'Missing query parameter'}), 400
    
    # Get video metadata
    metadata = get_video_metadata(query, username, password)
    return jsonify(metadata)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
