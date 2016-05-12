package com.khizerhasan.attendanceble;

import android.animation.Animator;
import android.animation.AnimatorListenerAdapter;
import android.annotation.TargetApi;
import android.bluetooth.le.AdvertiseCallback;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.content.IntentFilter;
import android.content.SharedPreferences;
import android.os.AsyncTask;
import android.os.Build;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;

import java.io.IOException;
import java.net.HttpURLConnection;

public class MainActivity extends AppCompatActivity {

    private BroadcastReceiver advertisingFailureReceiver;

    private View progressView;
    private String Name;
    private String SJSUID;

    private RegisterUserInfoTask registerUserInfoTask = null;

    public static final String ADVERTISING_FAILED = "com.khizerhasan.attendanceble.advertising_failed";

    public static final String ADVERTISING_FAILED_EXTRA_CODE = "failureCode";

    public static final int ADVERTISING_TIMED_OUT = 6;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        advertisingFailureReceiver = new BroadcastReceiver() {

            @Override
            public void onReceive(Context context, Intent intent) {
                int errorCode = intent.getIntExtra(MainActivity.ADVERTISING_FAILED_EXTRA_CODE, -1);


                String errorMessage = getString(R.string.error_prefix);
                switch (errorCode) {
                    case AdvertiseCallback.ADVERTISE_FAILED_ALREADY_STARTED:
                        errorMessage += " " + getString(R.string.error_already_started);
                        break;
                    case AdvertiseCallback.ADVERTISE_FAILED_DATA_TOO_LARGE:
                        errorMessage += " " + getString(R.string.error_too_large);
                        break;
                    case AdvertiseCallback.ADVERTISE_FAILED_FEATURE_UNSUPPORTED:
                        errorMessage += " " + getString(R.string.error_unsupported);
                        break;
                    case AdvertiseCallback.ADVERTISE_FAILED_INTERNAL_ERROR:
                        errorMessage += " " + getString(R.string.error_internal);
                        break;
                    case AdvertiseCallback.ADVERTISE_FAILED_TOO_MANY_ADVERTISERS:
                        errorMessage += " " + getString(R.string.error_too_many);
                        break;
                    case MainActivity.ADVERTISING_TIMED_OUT:
                        errorMessage = " " + getString(R.string.advertising_timedout);
                        break;
                    default:
                        errorMessage += " " + getString(R.string.error_unknown);
                }

                Toast.makeText(getBaseContext(), errorMessage, Toast.LENGTH_LONG).show();
            }
        };




        progressView = findViewById(R.id.register_user_progress);
        final EditText firstName = (EditText) findViewById(R.id.user_firstName);
        final EditText lastName = (EditText) findViewById(R.id.user_lastName);
        final EditText sjsuId = (EditText) findViewById(R.id.user_sjsuId);
        Button registerButton = (Button) findViewById(R.id.register_button);

        registerButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                SharedPreferences pref = getSharedPreferences("UserDetails",0);
                SharedPreferences.Editor edit = pref.edit();
                Name = firstName.getText().toString()+"."+lastName.getText().toString();
                SJSUID = sjsuId.getText().toString();
                edit.putString("username",Name);
                edit.putString("sjsuid",SJSUID);
                edit.commit();

                showProgress(true);

                registerUserInfoTask = new RegisterUserInfoTask(getApplicationContext());
                registerUserInfoTask.execute((Void) null);
                finish();


            }
        });


    }

    @Override
    protected void onResume() {
        super.onResume();

        SharedPreferences pref = getSharedPreferences("UserDetails",0);
        String id = pref.getString("sjsuid","");
        if( id.equals("")){

        }
        else{
            Intent i = new Intent(MainActivity.this,Status.class);
            startActivity(i);
            finish();
        }
        IntentFilter failureFilter = new IntentFilter(MainActivity.ADVERTISING_FAILED);
        getBaseContext().registerReceiver(advertisingFailureReceiver, failureFilter);

    }

    protected void onPause() {
        super.onPause();
        getBaseContext().unregisterReceiver(advertisingFailureReceiver);
    }

    public class RegisterUserInfoTask extends AsyncTask<Void, Void, Boolean> {
        private Context mainContext;
        RegisterUserInfoTask(Context context) {
            mainContext = context;
        }

        @Override
        protected Boolean doInBackground(Void... params) {
            String url = "http://52.37.72.212:3005/"+SJSUID+"/"+ Name;
                HttpConnectionHelper connectionHelper;
                int returnCode;
               try {
                    connectionHelper = new HttpConnectionHelper(url, "PUT", HttpConnectionHelper.DEFAULT_CONNECT_TIME_OUT);
                    //connectionHelper.setRequestProperty("Content-type", "application/json");
                returnCode = connectionHelper.request_noOutput();
            } catch (IOException e) {
                return false;
            }
            return HttpURLConnection.HTTP_NO_CONTENT == returnCode;
            //return true;
        }

        @Override
        protected void onPostExecute(final Boolean success) {
            registerUserInfoTask = null;
            showProgress(false);
            if (success) {
                try {
                    Toast.makeText(mainContext,"Advertising Service Started",Toast.LENGTH_SHORT).show();
                    Intent in = new Intent(mainContext,AsyncTask.Status.class);
                    startActivity(in);

                   // Intent i = new Intent(MainActivity.this,TakeAction.class);
                    //startActivity(i);
                } catch (Exception e) {
                    Log.d("Exception",e.getMessage());
                }
            } else {
                //TODO: Show no content?
            }
        }

        @Override
        protected void onCancelled() {
            registerUserInfoTask = null;
            showProgress(false);
        }
    }

    @TargetApi(Build.VERSION_CODES.HONEYCOMB_MR2)
    private void showProgress(final boolean show) {
        // On Honeycomb MR2 we have the ViewPropertyAnimator APIs, which allow
        // for very easy animations. If available, use these APIs to fade-in
        // the progress spinner.
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.HONEYCOMB_MR2) {
            int shortAnimTime = getResources().getInteger(android.R.integer.config_shortAnimTime);

            progressView.setVisibility(show ? View.VISIBLE : View.GONE);
            progressView.animate().setDuration(shortAnimTime).alpha(
                    show ? 1 : 0).setListener(new AnimatorListenerAdapter() {
                @Override
                public void onAnimationEnd(Animator animation) {
                    progressView.setVisibility(show ? View.VISIBLE : View.GONE);
                }
            });
        } else {
            // The ViewPropertyAnimator APIs are not available, so simply show
            // and hide the relevant UI components.
            progressView.setVisibility(show ? View.VISIBLE : View.GONE);
        }
    }

}
