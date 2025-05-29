import 'dart:io';
import 'package:shelf/shelf.dart';
import 'package:shelf/shelf_io.dart' as io;
import 'package:shelf_proxy/shelf_proxy.dart';

void main() async {
  final handler = proxyHandler('https://localhost:8008/',
      client: HttpClient()
        ..badCertificateCallback = (cert, host, port) => true);

  final pipeline =
      const Pipeline().addMiddleware(logRequests()).addHandler(handler);

  await io.serve(pipeline, 'localhost', 8080);
  print('Proxy server listening on port 8080');
}
