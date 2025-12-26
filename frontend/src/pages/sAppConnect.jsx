import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import QRCode from 'react-qr-code';

function WhatsAppConnect() {
  const [qrCode, setQrCode] = useState('');
  const [status, setStatus] = useState('disconnected');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchQR = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/get-whatsapp-qr`);
        if (response.ok) {
          const data = await response.json();
          setStatus(data.status);
          setQrCode(data.qr);
        }
      } catch (err) {
        console.error("Error fetching QR:", err);
      }
    };

    // Check immediately, then every 2 seconds
    fetchQR();
    const interval = setInterval(fetchQR, 2000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="container">
      <h2>Connect WhatsApp</h2>
      
      {status === 'connected' ? (
        <div>
          <h3>Connected!</h3>
          <p>Your phone is linked successfully.</p>
        </div>
      ) : (
        <div>
          {qrCode ? (
            <div>
              <QRCode value={qrCode} />
            </div>
          ) : (
            <p>Loading QR Code...</p>
          )}
        </div>
      )}

      <br /><br />
      <button onClick={() => navigate("/")}>
        Home
      </button>
    </div>
  );
}

export default WhatsAppConnect;