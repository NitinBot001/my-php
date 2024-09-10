require 'sinatra'
require 'json'
require 'open3'

# Helper method to extract video details using yt-dlp
def fetch_video_details(video_url)
  # Command to extract video details using yt-dlp
  command = "yt-dlp -J --no-playlist #{video_url}"

  stdout, stderr, status = Open3.capture3(command)

  if status.success?
    video_data = JSON.parse(stdout)

    title = video_data['title']
    thumbnail_url = video_data['thumbnail']
    audio_link = video_data['url']  # Direct audio play link

    { title: title, thumbnail_url: thumbnail_url, audio_link: audio_link }
  else
    { error: stderr }
  end
end

# Route to handle the request and fetch video details
get '/video_info' do
  content_type :json

  video_url = params[:video_url]

  if video_url.nil? || video_url.empty?
    halt 400, { error: "video_url parameter is required" }.to_json
  end

  video_details = fetch_video_details(video_url)
  video_details.to_json
end

# Start the Sinatra server
run! if app_file == $0
