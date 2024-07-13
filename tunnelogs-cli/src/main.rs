use clap::Arg;
use std::io::BufRead;
use std::sync::mpsc;
use std::thread;
use websocket::{ClientBuilder, OwnedMessage};

fn main() {
    let matches = clap::Command::new("tunnelogs")
        .version("0.1")
        .about("Tunnel logs sends logs to server to specific lobby")
        .arg(Arg::new("lobby").required(true))
        .get_matches();

    let lobby_name = matches
        .get_one::<String>("lobby")
        .expect("lobby is required");

    let (stdin_sender, stdin_receiver) = mpsc::channel();
    let url = format!("ws://localhost:8080/connect/{}/server", lobby_name);
    let (close_sender, close_receiver) = mpsc::channel::<()>();

    thread::spawn(move || {
        let stdin = std::io::stdin();
        let mut handle = stdin.lock();

        // Read input line by line
        let mut message = String::new();
        while handle.read_line(&mut message).unwrap() > 0 {
            println!("Received message {}", message);
            // Clear the buffer for the next line
            stdin_sender.send(message.trim().to_string()).unwrap();
            message.clear();
        }
    });

    thread::spawn(move || {
        let client = match ClientBuilder::new(&url).unwrap().connect_insecure() {
            Ok(client) => client,
            Err(err) => {
                println!("failed to connect to ws server {err}");
                close_sender.send(()).unwrap();
                return;
            }
        };

        let (_receiver, mut sender) = client.split().unwrap();

        // Send messages
        while let Ok(message) = stdin_receiver.recv() {
            if let Err(err) = sender.send_message(&OwnedMessage::Text(message)) {
                println!("failed to send message {err}");
                close_sender.send(()).unwrap();
                return;
            }
        }

        close_sender.send(()).unwrap();
    });

    close_receiver.recv().unwrap();
}
