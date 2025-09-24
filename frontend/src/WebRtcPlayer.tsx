/*
 * @Description:
 * @Version: 1.0
 * @Author: Knight
 * @Date: 2024-01-11 15:54:42
 * @LastEditors: Knight
 * @LastEditTime: 2025-04-17 15:39:12
 */
import { FC, useEffect, useRef, useState } from "react";
import { Play } from "../wailsjs/go/main/App";
import Loading from "./Loading";

export interface WebRtcPlayerProps {
  url?: string;
}

const WebRtcPlayer: FC<WebRtcPlayerProps> = (props) => {
  const { url } = props;
  const videoRef = useRef<HTMLVideoElement>(null!);
  const rtcPeerCon = useRef<RTCPeerConnection | null>(null);
  const isLoadData = useRef(false);
  const [loading, setLoading] = useState<boolean>(false);

  const onDoubleClick = (): void => {
    if (isLoadData.current) {
      videoRef.current.requestPictureInPicture();
    }
  };

  const onLoadedData = (): void => {
    isLoadData.current = true;
    videoRef.current.play()?.finally(() => setLoading(false));
  };

  const open = async (): Promise<void> => {
    if (!url) {
      return;
    }

    const pc = (rtcPeerCon.current = new RTCPeerConnection());
    const stream = new MediaStream();

    pc.onnegotiationneeded = async (): Promise<void> => {
      setLoading(true);
      videoRef.current.load();
      videoRef.current.src = "";
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);
      if (!pc.localDescription?.sdp) {
        return;
      }

      const res = await Play(url, btoa(pc.localDescription.sdp));

      if (!res || !pc.setRemoteDescription || pc.signalingState === "closed") {
        return;
      }

      pc.setRemoteDescription(
        new RTCSessionDescription({
          type: "answer",
          sdp: atob(res),
        })
      );
    };

    pc.ontrack = (event): void => {
      stream.addTrack(event.track);
      videoRef.current.srcObject = stream;
    };

    pc.onsignalingstatechange = (): void => {
      console.log(`onsignalingstatechange:${pc.signalingState}(${url})`);
    };

    pc.onconnectionstatechange = (): void => {
      console.log(`connectionState:${pc.connectionState}(${url})`);
    };

    pc.onicegatheringstatechange = (): void => {
      console.log(`iceGatheringState:${pc.iceGatheringState}(${url})`);
    };

    pc.oniceconnectionstatechange = (): void => {
      console.log(`iceConnectionState:${pc.iceConnectionState}(${url})`);

      if (pc.iceConnectionState === "disconnected") {
        // 视频出错后销毁重建
        close();
        open();
      }
    };

    pc.addTransceiver("video", { direction: "sendrecv" });
  };

  const close = (): void => {
    videoRef.current?.pause();
    rtcPeerCon.current?.close();
    rtcPeerCon.current = null;
  };

  useEffect(() => {
    open();

    return () => {
      close();
    };
  }, [url]);

  return (
    <div className="relative w-full bg-black rounded overflow-hidden">
      {loading && (
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 text-white text-2xl">
          <Loading />
        </div>
      )}

      <video
        muted
        autoPlay={false}
        controls={false}
        ref={videoRef}
        className="w-full aspect-video align-middle"
        onDoubleClick={onDoubleClick}
        onLoadedData={onLoadedData}
      />
    </div>
  );
};

export default WebRtcPlayer;
