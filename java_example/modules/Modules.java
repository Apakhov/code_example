import java.util.*;
import java.util.function.Predicate;

public class Modules {
    public static void main(String[] argv){
        Scanner in = new Scanner(System.in);
        StringBuilder input = new StringBuilder();
        while (in.hasNext())
            input.append(in.next());

        OrientedGraph funcGraph;
        try{
        funcGraph = (new Parser(new Tokens(input.toString()))).getFuncGraph();
        }catch (Exception e){
            System.out.println("error");
            return;
        }

        int compAm = funcGraph.getStrongComponent();
        int[] compCount = new int[compAm];
        for(Vertex v : funcGraph.getVertices()){
            compCount[v.getComponent()]++;
        }

        int counter = funcGraph.getVertices().size();
        for(int x : compCount){
            counter -= (x - 1);
        }
        System.out.println(counter);
    }
}

class Token{
    public enum TYPE{
        IDENTIFIER,
        NUMBER,
        ADDOPER,
        MULOPER,
        LPARENTHESIS,
        RPARENTHESIS,
        COLON,
        FUNCANN,
        QUESTION,
        SEMICOLON,
        COMMA,
        COMPARATOR
    }
    
    public final TYPE type;
    public final String content;

    Token(TYPE type, String content){
        this.type = type;
        this.content = content;
    }

    @Override
    public String toString() {
        return "[" + type + "]" + content;
    }
}

class Tokens{
    private ArrayList<Token> tokens;
    private int curToken;

    private int skip(int i, String text, Predicate<Integer> r){
        while (i < text.length() && r.test(text.codePointAt(i))) i++;

        return i;
    }

    Tokens(String text){
        tokens = new ArrayList<>();
        int len = text.length();
        int i = 0;
        while (i < len){
            i = skip(i, text, Character::isWhitespace);
            if(i == len) break;
            int begin = i;
            Token.TYPE type;
            int curSymb = text.codePointAt(i);
            if(curSymb == '+' || curSymb == '-'){
                type = Token.TYPE.ADDOPER;
                i++;
            }else if(curSymb == '*' || curSymb == '/'){
                type = Token.TYPE.MULOPER;
                i++;
            }else if(curSymb == '('){
                type = Token.TYPE.LPARENTHESIS;
                i++;
            }else if(curSymb == ')'){
                type = Token.TYPE.RPARENTHESIS;
                i++;
            }else if(curSymb == ';'){
                type = Token.TYPE.SEMICOLON;
                i++;
            }else if(curSymb == '?'){
                type = Token.TYPE.QUESTION;
                i++;
            }else if(curSymb == ','){
                type = Token.TYPE.COMMA;
                i++;
            }else if(curSymb == '='){
                type = Token.TYPE.COMPARATOR;
                i++;
            }else if(curSymb == '>'){
                type = Token.TYPE.COMPARATOR;
                i++;
                if(i != len && text.codePointAt(i) == '=')
                    i++;
            }else if(curSymb == '<'){
                type = Token.TYPE.COMPARATOR;
                i++;
                if(i != len && (text.codePointAt(i) == '=' || text.codePointAt(i) == '>'))
                    i++;
            }else if(curSymb == ':'){
                if(i == len - 1 || text.codePointAt(i + 1) != '='){
                    type = Token.TYPE.COLON;
                    i++;
                }else{
                    type = Token.TYPE.FUNCANN;
                    i+=2;
                }
            }else if(Character.isDigit(curSymb)){
                i = skip(i, text, Character::isDigit);
                if(i < len && Character.isAlphabetic(text.codePointAt(i))){
                    throw new RuntimeException("delimiter expected");
                }
                type = Token.TYPE.NUMBER;
            }else if(Character.isAlphabetic(curSymb)){
                i = skip(i, text, Character::isLetterOrDigit);
                type = Token.TYPE.IDENTIFIER;
            }else {
                throw new RuntimeException("wrong symbol: " + i + " : !" + (char)text.codePointAt(i) + "!");
            }


            tokens.add(new Token(type, text.substring(begin, i)));
        }
        curToken = 0;
    }

    public boolean hasNext(){
        return curToken < tokens.size();
    }

    public Token peek(){
        if(!hasNext()){
            throw new RuntimeException("End of tokens");
        }
        return tokens.get(curToken);
    }

    public Token next(){
        if(!hasNext()){
            throw new RuntimeException("End of tokens");
        }
        return tokens.get(curToken++);
    }

    public void reload(){
        curToken = 0;
    }

    @Override
    public String toString() {
        StringBuilder s = new StringBuilder();
        tokens.forEach(x -> {s.append(x); s.append("\n");});
        return s.toString();
    }
}

class Parser{
    private static final int UNFOUND = 0x00;
    private static final int FOUNDDECLARATION = 0x01;
    private static final int FOUNDUSAGE =  0x10;
    private static final int FOUNDDECLANDUSE = 0x11;
    private Tokens tokens;

    Parser(Tokens tokens){
            this.tokens = tokens;
    }

    public  OrientedGraph getFuncGraph()throws Exception {
        OrientedGraph orGr = new OrientedGraph();
        TreeMap<String, Integer> foundFuncs = new TreeMap<>();
        TreeMap<String, TreeSet<String>> usageMap = new TreeMap<>();
        TreeMap<String, Integer> argMap = new TreeMap<>();

        parseProgOG(foundFuncs, usageMap, argMap);

        if(foundFuncs.containsValue(FOUNDUSAGE))
            throw new Exception("undeclared function was found");

        TreeMap<String, Vertex> graphBuilder = new TreeMap<>();
        usageMap.forEach((x,y)->{
            graphBuilder.put(x, new Vertex(x));
        });
        usageMap.forEach((x,y)->{
            orGr.addVertex(graphBuilder.get(x));
            y.forEach(z -> graphBuilder.get(x).addEdge(graphBuilder.get(z)));
        });

        return orGr;
    }

    private Token mustBeOG(Token.TYPE type)throws Exception{
        if(!tokens.hasNext()){
            throw new RuntimeException("expected: " + type + ", but nothing found");
        }
        Token token = tokens.next();
        if(token.type != type){
            throw new RuntimeException("token: " + token.content + ", expected: " + type + ", but found: " + token.type);
        }
        return token;
    }

    private void parseProgOG(TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap)throws Exception{
        while (tokens.hasNext()){
            parseFuncOG(foundFuncs, usageMap, argMap);
        }
        tokens.reload();
    }

    private void parseFuncOG(TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap)throws Exception{
        Token curToken = mustBeOG(Token.TYPE.IDENTIFIER);
        Integer val = foundFuncs.get(curToken.content);
        if(val == null){
            foundFuncs.put(curToken.content, FOUNDDECLARATION);
        }else {
            foundFuncs.put(curToken.content, FOUNDDECLARATION | val);
        }
        mustBeOG(Token.TYPE.LPARENTHESIS);
        TreeSet<String> idents = parseFormArgOG();

        if(usageMap.containsKey(curToken.content)){
            if(argMap.get(curToken.content) != idents.size())
                throw new Exception("wrong amount of parameters in " + curToken.content);
        }else {
            usageMap.put(curToken.content, new TreeSet<>());
            argMap.put(curToken.content, idents.size());
        }

        mustBeOG(Token.TYPE.RPARENTHESIS);
        mustBeOG(Token.TYPE.FUNCANN);
        parseExprOG(curToken.content, foundFuncs, usageMap, argMap, idents);
        mustBeOG(Token.TYPE.SEMICOLON);
    }

    private TreeSet<String> parseFormArgOG()throws Exception {
        TreeSet<String> foundTokens = new TreeSet<>();
        while (tokens.hasNext() && tokens.peek().type == Token.TYPE.IDENTIFIER){
            foundTokens.add(tokens.next().content);
            if(tokens.hasNext() && tokens.peek().type == Token.TYPE.COMMA){
                tokens.next();
            if(!tokens.hasNext() || tokens.peek().type != Token.TYPE.IDENTIFIER)
                throw new  Exception("msg");
            }
        }
        return foundTokens;
    }

    private void parseExprOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        parseCompExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
        if(tokens.hasNext() && tokens.peek().type == Token.TYPE.QUESTION){
            tokens.next();
            parseCompExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
            mustBeOG(Token.TYPE.COLON);
            parseExprOG(curFunc, foundFuncs, usageMap, argMap, idents);
        }
    }

    private void parseCompExpOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        parseArBegExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
        if(tokens.hasNext() && tokens.peek().type == Token.TYPE.COMPARATOR){
            tokens.next();
            parseArBegExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
        }
    }

    private void parseArBegExpOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        parseTerBegExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
        parseArEndExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
    }

    private void parseArEndExpOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        if(tokens.hasNext() && tokens.peek().type == Token.TYPE.ADDOPER){
            tokens.next();
            parseTerBegExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
            parseArEndExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
        }
    }

    private void parseTerBegExpOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        parseFactorOG(curFunc, foundFuncs, usageMap, argMap, idents);
        parseTerEndExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
    }

    private void parseTerEndExpOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        if(tokens.hasNext() && tokens.peek().type == Token.TYPE.MULOPER){
            tokens.next();
            parseFactorOG(curFunc, foundFuncs, usageMap, argMap, idents);
            parseTerEndExpOG(curFunc, foundFuncs, usageMap, argMap, idents);
        }
    }

    private void parseFactorOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        if(!tokens.hasNext())
            throw new Exception("expected expression, but found nothing");
        Token curToken = tokens.next();
        switch (curToken.type){
            case NUMBER:
                return;
            case ADDOPER:
                if(curToken.content.equals("-"))
                    parseFactorOG(curFunc, foundFuncs, usageMap, argMap, idents);
                return;
            case LPARENTHESIS:
                parseExprOG(curFunc, foundFuncs, usageMap, argMap, idents);
                mustBeOG(Token.TYPE.RPARENTHESIS);
                return;
            case IDENTIFIER:
                if(!tokens.hasNext() || tokens.peek().type != Token.TYPE.LPARENTHESIS) {
                    if (idents.contains(curToken.content))
                        return;
                    else
                        throw new Exception("undefined variable");
                }else {
                    tokens.next();
                    Integer val = foundFuncs.get(curToken.content);
                    if(val == null){
                        foundFuncs.put(curToken.content, FOUNDUSAGE);
                    }else {
                        foundFuncs.put(curToken.content, FOUNDUSAGE | val);
                    }
                    int args = parseActArgOG(curFunc, foundFuncs, usageMap, argMap, idents);
                    if(usageMap.containsKey(curToken.content)){
                        if(argMap.get(curToken.content) != args)
                            throw new Exception("wrong amount of parameters in " + curToken.content + "expected: " + argMap.get(curToken.content) + "got: " + args + "\n in defenition of " + curFunc);
                    }else {
                        usageMap.put(curToken.content, new TreeSet<>());
                        argMap.put(curToken.content, args);
                    }
                    usageMap.get(curFunc).add(curToken.content);
                    mustBeOG(Token.TYPE.RPARENTHESIS);
                    return;
                }
            default:
                throw new Exception("expected expression, but found: " + curToken.content);

        }
    }

    private int parseActArgOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        if(tokens.peek().type != Token.TYPE.RPARENTHESIS){
            return parseExpListOG(curFunc, foundFuncs, usageMap, argMap, idents);
        }
        return 0;
    }

    private int parseExpListOG(String curFunc, TreeMap<String, Integer> foundFuncs, TreeMap<String, TreeSet<String>> usageMap, TreeMap<String, Integer> argMap, TreeSet<String> idents)throws Exception{
        parseExprOG(curFunc, foundFuncs, usageMap, argMap, idents);
        if(tokens.peek().type == Token.TYPE.COMMA){
            tokens.next();
            return 1 + parseExpListOG(curFunc, foundFuncs, usageMap, argMap, idents);
        }
        return 1;
    }


}

class Vertex{
    private ArrayList<Vertex> edges;
    private String name;
    private int time;
    private int component;
    private int low;

    Vertex(String name){
        this.name = name;
        edges = new ArrayList<>();
        time = -1;
        component = -1;
        low = -1;
    }

    public void addEdge(Vertex v){
        edges.add(v);
    }

    public ArrayList<Vertex> getEdges() {
        return edges;
    }

    public String getName() {
        return name;
    }

    public void setTime(int time) {
        this.time = time;
    }

    public int getTime() {
        return time;
    }

    public void setComponent(int component) {
        this.component = component;
    }

    public int getComponent() {
        return component;
    }

    public void setLow(int low) {
        this.low = low;
    }

    public int getLow() {
        return low;
    }

    @Override
    public String toString() {
        return name;
    }
}


class OrientedGraph{
    private ArrayList<Vertex> vertices;
    OrientedGraph(){
        vertices = new ArrayList<>();
    }

    public void addVertex(Vertex v){
        vertices.add(v);
    }

    public ArrayList<Vertex> getVertices() {
        return vertices;
    }

    private int time;
    private int component;

    public int getStrongComponent(){
        time = 1;
        component = 0;
        for(Vertex v: vertices){
            v.setComponent(-1);
            v.setTime(0);
        }
        Stack<Vertex> stack = new Stack<>();
        for(Vertex v: vertices){
            if (v.getTime() == 0)
                visitVertexTarjan(v, stack);
        }
        return component;
    }

    private void visitVertexTarjan(Vertex v, Stack<Vertex> stack){
        v.setTime(time);
        v.setLow(time);
        time++;
        stack.push(v);
        for(Vertex u: v.getEdges()){
            if(u.getTime() == 0)
                visitVertexTarjan(u, stack);
            if(u.getComponent() == -1 && v.getLow() > u.getLow())
                v.setLow(u.getLow());
        }
        if(v.getLow() == v.getTime()){
            Vertex w;
            do {
                w = stack.pop();
                w.setComponent(component);
            }while (w != v);
            component++;
        }
    }

    public String getGraphViz(){
        StringBuilder out = new StringBuilder();
        out.append("digraph{");
        for(Vertex v : vertices){
            out.append(v.getName());
            out.append("\n");
            for(Vertex x : v.getEdges()){
                out.append(v.getName());
                out.append("->");
                out.append(x.getName());
                out.append("\n");
            }
        }
        out.append("}");

        return out.toString();
    }

    @Override
    public String toString() {
        StringBuilder out = new StringBuilder();
        for(Vertex v : vertices){
            out.append("[");
            out.append(v.getName());
            out.append("]");
            for(Vertex x : v.getEdges()){
                out.append(" ");
                out.append(x.getName());
            }
            out.append("\n");
        }

        return out.toString();
    }
}