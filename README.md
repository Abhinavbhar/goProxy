# Proxy Server with Auth & Chrome Extension

### Why use this proxy?
Sometimes, the internet isn’t as open as it should be. You may face restrictions like:
- **Geo-restricted content** – websites or videos only available in certain regions.  
- **Public Wi-Fi limitations** – schools, offices, airports, or cafes often block social media, streaming sites, or other services.  
- **Network-level censorship** – some ISPs restrict or throttle access to specific platforms.  

This project provides a simple way around such restrictions by routing your traffic through a **forward proxy** with authentication and easy setup using a **Chrome extension**.

---

### How it works
- **Forward Proxy Server**  
  A lightweight proxy built to forward HTTP(S) requests. It acts as the middle layer between your browser and the internet, so websites see requests coming from the proxy instead of your actual network.  

- **Backend Authentication**  
  A minimal backend ensures that only authorized users can use the proxy. This prevents open abuse and keeps your server secure.  

- **Chrome Extension**  
  A small extension automatically configures your browser to use the proxy once you’re logged in. This saves you from manual proxy settings and makes the experience seamless.  

---

### Setup Instructions

#### 1. Clone the repository
```bash
git clone https://github.com/your-username/proxy-server.git
cd proxy-server

2. Start MongoDB

Make sure MongoDB is running locally at:

mongodb://localhost:27017

3. Run the Proxy Server

cd proxy
go run .

4. Run the Backend

cd backend
go run main.go

5. Load the Chrome Extension

    Open Chrome and go to chrome://extensions/

    Enable Developer mode (top-right corner)

    Click Load unpacked

    Select the extension/ folder from this project

