from flask import Flask, request, abort
import pychromecast
from dataclasses import dataclass, asdict
import json
from uuid import UUID

app = Flask(__name__)

cast_sessions = {}

# Get devices

@dataclass
class CastDevice:
    name: str
    model: str
    type: str
    uuid: str

@app.route("/get_devices", methods=["GET"])
def get_devices():
    
    services, browser = pychromecast.discovery.discover_chromecasts()
    pychromecast.discovery.stop_discovery(browser)

    cast_list = []

    for s in services:

        cast_list.append(
            CastDevice(
                name=s.friendly_name,
                model=s.model_name,
                type=s.cast_type,
                uuid=str(s.uuid)
            )
        )
    
    return json.dumps([asdict(x) for x in cast_list], indent=4)

# Set volume

@dataclass(kw_only=True)
class SetVolumeQuery:
    uuid: str
    volume: float

@app.route("/set_volume", methods=["POST"])
def set_volume():
    query = SetVolumeQuery(**request.get_json())

    if not query.uuid in cast_sessions:
        chromecasts, browser = pychromecast.get_listed_chromecasts(uuids=[UUID(query.uuid)])
    
        cast_sessions[query.uuid] = chromecasts[0]
        cast_sessions[query.uuid].wait()

        pychromecast.discovery.stop_discovery(browser)

    cast_sessions[query.uuid].set_volume(query.volume)

    return ""

# Play media

@dataclass(kw_only=True)
class PlayMediaQuery:
    url: str
    mimetype: str
    uuid: str

@app.route("/play_media", methods=["POST"])
def play_media():
    query = PlayMediaQuery(**request.get_json())

    if not query.uuid in cast_sessions:
        chromecasts, browser = pychromecast.get_listed_chromecasts(uuids=[UUID(query.uuid)])
    
        cast_sessions[query.uuid] = chromecasts[0]
        cast_sessions[query.uuid].wait()

        pychromecast.discovery.stop_discovery(browser)

    controller = cast_sessions[query.uuid].media_controller
    controller.play_media(query.url, query.mimetype)
    controller.block_until_active()

    return ""

# Media control

@dataclass(kw_only=True)
class ControlQuery:
    uuid: str
    state: str

@app.route("/control", methods=["POST"])
def control():
    query = ControlQuery(**request.get_json())

    if not query.uuid in cast_sessions:
        chromecasts, browser = pychromecast.get_listed_chromecasts(uuids=[UUID(query.uuid)])
    
        cast_sessions[query.uuid] = chromecasts[0]
        cast_sessions[query.uuid].wait()

        pychromecast.discovery.stop_discovery(browser)

    controller = cast_sessions[query.uuid].media_controller
    controller.block_until_active()

    match query.state:
        case "play":
            controller.play()

        case "pause":
            controller.pause()
            
        case "stop":
            controller.stop()
        case _:
            return abort(400)
    
    return ""

# Device status

@dataclass
class DeviceStatus:
    volume: float
    current_time: float

@app.route("/get_status/<uuid>")
def status(uuid=None):
    
    if uuid == None or uuid == "":
        return abort(400)
    
    # This has to run again otherwise the MediaStatus remains the same

    chromecasts, browser = pychromecast.get_listed_chromecasts(uuids=[UUID(uuid)])

    cast_sessions[uuid] = chromecasts[0]
    cast_sessions[uuid].wait()

    pychromecast.discovery.stop_discovery(browser)

    controller = cast_sessions[uuid].media_controller
    controller.block_until_active()

    return json.dumps(asdict(DeviceStatus(
        volume=controller.status.volume_level,
        current_time=controller.status.current_time
    )))

if __name__ == "__main__":
    app.run(host="0.0.0.0", port="8080")