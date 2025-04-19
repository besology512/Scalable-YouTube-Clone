document.getElementById("play").onclick = () => {
    const name = document.getElementById("videoName").value.trim();
    if (!name) {
      return alert("Please enter a stream name");
    }
  
    // **Point directly at your streaming‑service** on port 8083
    const url = `http://localhost:8083/hls/${encodeURIComponent(name)}/index.m3u8`;
    const video = document.getElementById("player");
  
    if (Hls.isSupported()) {
      const hls = new Hls();
      hls.loadSource(url);
      hls.attachMedia(video);
      hls.on(Hls.Events.MANIFEST_PARSED, () => video.play());
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
      video.src = url;
      video.addEventListener("loadedmetadata", () => video.play());
    } else {
      alert("HLS not supported in this browser");
    }
  };
  