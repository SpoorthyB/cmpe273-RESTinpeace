package com.khizerhasan.attendanceble;

import android.content.Intent;
import android.content.SharedPreferences;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.TextView;

public class Status extends AppCompatActivity {

    private int serviceRunning = 0;
    private Button startServiceButton;
    private Button stopServiceButton;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_status);

        Intent i = new Intent(getBaseContext(),Advertise.class);
        startService(i);
        serviceRunning = 1;
        SharedPreferences pref = getSharedPreferences("UserDetails",0);
        String id = pref.getString("sjsuid","");

        TextView sjsuid = (TextView) findViewById(R.id.sjsuId);
        sjsuid.setText(id);


         startServiceButton= (Button) findViewById(R.id.start_button);
         stopServiceButton= (Button) findViewById(R.id.stop_button);

        startServiceButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(getBaseContext(),Advertise.class);
                startService(i);
            }
        });

        stopServiceButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent i = new Intent(getBaseContext(),Advertise.class);
                stopService(i);
            }
        });
    }

    @Override
    protected void onResume() {
        super.onResume();
        if(serviceRunning ==1){
            stopServiceButton.setEnabled(true);
            startServiceButton.setEnabled(false);
        }
        else{
            stopServiceButton.setEnabled(false);
            startServiceButton.setEnabled(true);
        }
    }
}
