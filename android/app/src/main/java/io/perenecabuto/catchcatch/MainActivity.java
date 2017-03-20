package io.perenecabuto.catchcatch;

import android.app.Activity;
import android.content.Context;
import android.location.Location;
import android.location.LocationManager;
import android.os.Bundle;
import android.support.annotation.NonNull;
import android.support.v4.app.ActivityCompat;
import android.util.Log;
import android.widget.Toast;

import com.google.android.gms.maps.CameraUpdateFactory;
import com.google.android.gms.maps.GoogleMap;
import com.google.android.gms.maps.MapFragment;
import com.google.android.gms.maps.model.LatLng;
import com.google.android.gms.maps.model.Marker;
import com.google.android.gms.maps.model.MarkerOptions;

import org.json.JSONException;

import java.net.URISyntaxException;
import java.util.HashMap;
import java.util.List;

import io.socket.client.IO;
import io.socket.client.Socket;

import static android.Manifest.permission.ACCESS_FINE_LOCATION;
import static android.location.LocationManager.NETWORK_PROVIDER;
import static android.support.v4.content.PermissionChecker.PERMISSION_GRANTED;


public class MainActivity extends Activity implements ConnectionManager.EventCallback {

    private static final String TAG = MainActivity.class.getSimpleName();
    private static final int LOCATION_PERMISSION_REQUEST_CODE = (int) (Math.random() * 10000);
    private String provider = NETWORK_PROVIDER;

    private MapFragment mapFragment;
    private GoogleMap map;
    private ConnectionManager manager;
    private HashMap<String, Marker> markers = new HashMap<>();
    private Player player = new Player("", 0, 0);

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        mapFragment = (MapFragment) getFragmentManager().findFragmentById(R.id.activity_main_map);
        mapFragment.getMapAsync((GoogleMap m) -> {
            map = m;
            // TODO OnMapCreate get features around and plot them
            m.setOnCameraMoveListener(() -> {
                // TODO OnCameraChange get features around and plot them
                // Log.d(TAG, "position: " + m.getCameraPosition().target + "zoom: " + m.getCameraPosition().zoom);
            });
        });

        setupConnection();

        if (checkCallingOrSelfPermission(ACCESS_FINE_LOCATION) == PERMISSION_GRANTED) {
            setupLocation();
        } else {
            requestPermission();
        }
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        manager.disconnect();
    }

    private void setupConnection() {
        try {
            Socket socket = IO.socket("http://192.168.23.102:5000");
            manager = new ConnectionManager(socket, this);
            manager.connect();
        } catch (URISyntaxException | ConnectionManager.NoConnectionException e) {
            e.printStackTrace();
            Log.e(TAG, e.getMessage(), e);
        }
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grants) {
        boolean permitted = requestCode == LOCATION_PERMISSION_REQUEST_CODE
            && grants.length > 0 && grants[0] == PERMISSION_GRANTED;

        if (permitted) {
            setupLocation();
        } else {
            requestPermission();
        }
    }

    private void requestPermission() {
        ActivityCompat.requestPermissions(this,
            new String[]{ACCESS_FINE_LOCATION}, LOCATION_PERMISSION_REQUEST_CODE);
    }

    @SuppressWarnings("MissingPermission")
    private void setupLocation() {
        LocationManager locationManager = (LocationManager) this.getSystemService(Context.LOCATION_SERVICE);
        locationManager.requestLocationUpdates(provider, 0, 0, new LocationUpdateListener((Location l) -> {
            Log.d(TAG, "location updated");
            try {
                manager.sendPosition(l);
            } catch (JSONException e) {
                e.printStackTrace();
                Toast.makeText(this, "Error to send position", Toast.LENGTH_SHORT).show();
            }
            player.updateLocation(l);
            showPlayerOnMap(player);
        }));
    }

    @Override
    public void onPlayerList(List<Player> players) {
        runOnUiThread(() -> {
            Log.d(TAG, "remote-player:list " + players);
            for (Player p : players) {
                showPlayerOnMap(p);
            }
        });
    }

    private void showPlayerOnMap(Player p) {
        Marker m = markers.get(p.getId());
        if (m == null) {
            Log.d(TAG, "--> add player" + p);
            m = map.addMarker(new MarkerOptions().position(p.getPoint()).title(p.getId()));
            markers.put(p.getId(), m);
        } else {
            Log.d(TAG, "--> update player" + p);
            m.setPosition(p.getPoint());
        }
        if (p.getId().equals(player.getId())) {
            map.moveCamera(CameraUpdateFactory.newLatLngZoom(p.getPoint(), 15));
        }
        m.setVisible(true);
        m.showInfoWindow();
    }

    @Override
    public void onRemotePlayerUpdate(Player p) {
        runOnUiThread(() -> {
            Log.d(TAG, "remote-player:updated " + p);
            showPlayerOnMap(p);
        });
    }

    @Override
    public void onRemoteNewPlayer(Player p) {
        runOnUiThread(() -> {
            Log.d(TAG, "remote-player:new " + p);
            showPlayerOnMap(p);
        });
    }

    @Override
    public void onRegistred(Player p) {
        this.player = p;
        runOnUiThread(() -> {
            Log.d(TAG, "player:registered " + p);
            showPlayerOnMap(p);
        });
    }

    @Override
    public void onRemotePlayerDestroy(Player p) {
        Marker m = markers.get(p.getId());
        markers.remove(p.getId());
        runOnUiThread(() -> m.remove());
    }
}
